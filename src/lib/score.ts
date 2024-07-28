export type Config = {
	language: "en" | "zh"
	teams: {
		name: string;
		admin: boolean | undefined;
	}[];
	resource_names: string[];
	gem_names: string[];
	lock_length_seconds: number;
};

export const blankConfig: Config = {
	teams: [],
	resource_names: [],
	gem_names: [],
	lock_length_seconds: 0,
}

export type Score = {
	resources: number[]
	gems: number[]
};

export type NoticeType = "pause" | "highlight" | "warning" | "message";
export type Notice = {
	timestamp: number
	id: number
	level: NoticeType
	message: string
	dismissed: boolean | undefined
	team_id: number | undefined
};

export const scoreBreakdown = (score: Score) => ({
	minResource: Math.min(...score.resources),
	gemType: score.gems.reduce((acc, num) => (num > 0 ? acc + 1 : acc), 0),
});

export const finalScore = (score: Score) => {
	const { minResource, gemType } = scoreBreakdown(score)
	return minResource * gemType
}
