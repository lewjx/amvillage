import { get, writable } from "svelte/store"

export type State = {
	config: Config
	balances: number[][]
	locks: (Lock | null)[]
	notice: string
	popup: string

	team: number
	username: string
}

export type Config = {
	lang: "en" | "zh"
	currencies: string[]
	gems: string[]
	teams: {
		name: string
		is_admin: boolean
	}[]
}

type Lock = {
	member: string
	millis_left: number
}

export type Cred = {
	username: string
	password: string
}

export const state = writable<State>({
	config: { lang: "en", currencies: [], gems: [], teams: [] },
	balances: [],
	locks: [],
	notice: "",
	team: 0,
	username: "",
	popup: "",
})

export const connected = writable(false)
export const error = writable("")
export const creds = writable<Cred|null>(null)

export const ws = writable<WebSocket>()

const tryLogin = () => {
	const isConnected = get(connected)
	const cred = get(creds)
	const websocket = get(ws)
	if (isConnected && cred && websocket) {
		websocket.send(`login ${btoa(cred.password)} ${cred.username}`)
	}
}

creds.subscribe(() => tryLogin())

export const connect = () => {
	const url = new URL((import.meta.env.VITE_BACKEND || "") + "/ws", location.href)
	url.protocol = location.protocol === "https:" ? "wss:" : "ws:"
	const websocket = new WebSocket(url)
	websocket.onopen = () => {
		connected.set(true)
		tryLogin()
	}
	websocket.onclose = () => {
		connected.set(false)
		console.warn("Websocket close")
		setTimeout(() => connect(), 1000)
	}
	websocket.onerror = err => {
		connected.set(false)
		console.warn("Websocket error:", err)
		setTimeout(() => connect(), 1000)
	}
	websocket.onmessage = (msg: MessageEvent) => {
		const data = msg.data as string
		if (data.indexOf(" ") === -1) {
			console.warn("expected space in event", data)
			return
		}
		const cmd = data.slice(0, data.indexOf(" "))
		const content = data.slice(data.indexOf(" "))
		switch (cmd) {
		case "state":
			console.log(JSON.parse(content))
			state.set(JSON.parse(content))
			break
		case "error":
			console.warn(content)
			error.set(content)
			break
		default:
			console.warn("unknown message", data)
		}
	}
	ws.set(websocket)
}
