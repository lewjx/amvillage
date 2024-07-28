export const disableScroll = () => {
	document.body.style.pointerEvents = "none";
	document.body.style.userSelect = "none";
	document.body.style.overflow = "hidden";
	document.body.style.height = "100%";
};
export const enableScroll = () => {
	document.body.style.pointerEvents = "auto";
	document.body.style.userSelect = "auto";
	document.body.style.overflow = "auto";
	document.body.style.height = "";
};
