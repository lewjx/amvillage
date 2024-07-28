import { type Session } from '$lib/login/login';
import type { Config } from '$lib/score';

export const isHoldingLock = (cfg: Config, session: Session) => {
	if (cfg.teams[session.team_id].admin) {
		// Admins implicitly hold the lock as they do not need it.
		return true;
	}
	const lockHolders = session.lock_holder[session.team_id].players
	return lockHolders.length > 0 && lockHolders[0].user_id === session.user_id;
}

export const isQueueing = (_: Config, session: Session) => {
	const lockHolders = session.lock_holder[session.team_id].players
	return lockHolders.some((x, i) => i !== 0 && x.user_id === session.user_id);
}

export const estimatedQueueingTime = (cfg: Config, session: Session) => {
	const lockInfo = session.lock_holder[session.team_id]
	const timeRemaining = lockInfo.seconds_remaining;
	const lockHolders = lockInfo.players
	// Calculate the wait time based on what we know.
	let location = lockHolders.findIndex((x) => x.user_id === session.user_id);
	if (location === -1) {
		// Assume that we will join right at the end of the queue if we are
		// not in the queue yet.
		location = lockHolders.length;
	}
	return timeRemaining + (location - 1) * cfg.lock_length_seconds;
}
