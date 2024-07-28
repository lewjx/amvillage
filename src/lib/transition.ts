export const maxHeight = (node: Element) => {
	const duration = 200
	const originalHeight = getComputedStyle(node).maxHeight;
	const targetHeight = originalHeight === 'none' ? node.scrollHeight : +originalHeight;
	return {
		delay: 0,
		duration,
		css: (t: number) => `max-height: ${t * targetHeight}px`
	};
};
