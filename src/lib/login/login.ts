import { get, writable, type Readable, type Writable } from 'svelte/store';
import { type Notice, type Score } from "$lib/score";

export type LockInfo = {
	seconds_remaining: number
	players: {
		nickname: string
		user_id: number
	}[]
}

export type Session = {
	ws?: WebSocket
	user_id: number
	team_id: number
	lock_holder: LockInfo[]
	score: Score[]
	notices: Notice[]
}

type StateMessage = { type: "state" } & Session;
type NoticeMessage = { type: "notice" } & Notice;

type WebSocketMessage = StateMessage | NoticeMessage;

export const createWebsocket = (name: string, secret: string): Promise<Readable<Session>> => {
	const store = writable<Session>({
		team_id: -1,
		user_id: -1,
		lock_holder: [],
		score: [],
		notices: [],
	});
	return new Promise((res, rej) => {
		keepWebsocketAlive(
			store, name, secret,
			() => res(store),
			(err: unknown) => {
				console.error(err)
				rej(err)
			},
		)
	});
}

const keepWebsocketAlive = (store: Writable<Session>, name: string, secret: string, res: () => void, rej: (err: unknown) => void) => {
	const url = new URL("api/ws", location.href)
	const ws = new WebSocket(url);
	let state: "disconnected" | "connected" | "established" = "disconnected"
	const retry = () => {
		switch (state) {
			case "disconnected":
				rej(new Error("Failed to connect to the server. Check your network connection."))
				return
			case "connected":
				// Do not retry if we did not successfully authenticate.
				rej(new Error("Failed to authenticate. Check your secret."))
				return
			case "established":
				break
		}
		store.update((original) => {
			original.ws = undefined
			setTimeout(() => keepWebsocketAlive(store, name, secret, res, rej), 1000)
			return original
		})
	}
	ws.onopen = () => {
		state = "connected"
		ws.send(JSON.stringify({
			type: "auth",
			nickname: name,
			secret,
		}))
	}
	ws.onclose = () => {
		console.warn("Websocket close")
		retry()
	}
	ws.onerror = err => {
		console.warn("Websocket error:", err)
		retry()
	}
	ws.onmessage = (msg: MessageEvent) => {
		if (state === "connected") {
			// First message, resolve the websocket.
			store.update((original) => {
				original.ws = ws
				return original
			})
			res()
			state = "established"
		}
		const latestInfo = JSON.parse(msg.data) as WebSocketMessage
		console.log(latestInfo)
		switch (latestInfo.type) {
			case "state":
				store.update((session) => {
					return {
						...session,
						...latestInfo,
					}
				})
				break
			case "notice":
				store.update((session) => {
					let notices = session.notices
					let updateIdx = session.notices.findIndex((x) => x.id === latestInfo.id)
					if (updateIdx !== -1) {
						session.notices.splice(updateIdx, 1)
					}
					return {
						...session,
						notices: [...notices, latestInfo],
					}
				})
				break
		}
		console.log(get(store))
	}
}
