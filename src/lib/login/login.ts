import { writable, type Readable, type Writable } from 'svelte/store';
import { type Message, type Score } from "$lib/score";

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
	messages: Message[]
}

export const createWebsocket = (name: string, secret: string): Promise<Readable<Session>> => {
	const store = writable<Session>({
		team_id: -1,
		user_id: -1,
		lock_holder: [],
		score: [],
		messages: [],
	});
	return new Promise((res, rej) => {
		keepWebsocketAlive(
			store, name, secret,
			() => res(store),
			rej,
		)
	});
}

const keepWebsocketAlive = (store: Writable<Session>, name: string, secret: string, res: () => void, rej: (err: unknown) => void) => {
	const ws = new WebSocket("/api/ws");
	let established = false
	const retry = () => {
		// Do not retry if we did not even receive one message.
		if (!established) {
			rej(new Error("Cannot establish WebSocket connection. Possible invalid secret?"))
			return
		}
		established = false
		store.update((original) => {
			original.ws = undefined
			setTimeout(() => keepWebsocketAlive(store, name, secret, res, rej), 1000)
			return original
		})
	}
	ws.onopen = () => {
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
		if (!established) {
			// First message, resolve the websocket.
			store.update((original) => {
				original.ws = ws
				return original
			})
			res()
			established = true
		}
		const latestInfo = JSON.parse(msg.data) as Session
		store.update((session) => {
			return {
				...session,
				...latestInfo,
			}
		})
		console.log(latestInfo)
	}
}
