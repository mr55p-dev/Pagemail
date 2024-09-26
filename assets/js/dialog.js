function openModal() {
    const elementId = this.getAttribute("data-id");
    document.getElementById(elementId).togglePopover();
}


function closeModalClickAway(event) {
    var rect = this.getBoundingClientRect();
    var isInDialog =
        rect.top <= event.clientY &&
        event.clientY <= rect.top + rect.height &&
        rect.left <= event.clientX &&
        event.clientX <= rect.left + rect.width;
    if (!isInDialog) {
        this.close();
    }
}
