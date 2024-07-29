import type { Config } from '$lib/score';
import type { Session } from '$lib/login/login';

export const isValid = (cfg: Config, session: Session, targetID: number, resourceCount: number[], gemCount: number[]) => {
	if (targetID < 0) {
		// If target is invalid, nothing is valid.
		return {
			resource: session.score[session.team_id].resources.map(_ => false),
			gem: session.score[session.team_id].gems.map(_ => false),
		}
	}
	const isAdmin = cfg.teams[session.team_id].admin;
	const ok = (ownArr: undefined | number[], targetArr: undefined | number[], requestedArr: number[]) => {
		return requestedArr.map((requested, i) => {
			const own = ownArr ? ownArr[i] : 0
			const target = targetArr ? targetArr[i] : 0
			if (typeof requested !== "number") return false
			switch (true) {
				case requested === 0:
					// OK if there are no transfers involved.
					return true
				case requested > 0:
					// OK if there is sufficient and we are transferring TO the target.
					// Also OK if we are admin.
					return own >= requested || isAdmin
				case requested < 0:
					// OK if we have permission AND the target has sufficient.
					return isAdmin && target >= -requested
			}
		})
	}
	return {
		resource: ok(session.score[session.team_id].resources, session.score[targetID].resources, resourceCount),
		gem: ok(session.score[session.team_id].gems, session.score[targetID].gems, gemCount),
	}
}
