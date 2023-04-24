function removeMarkedText() {
    const markedTextElements = document.querySelectorAll("mark");
    markedTextElements.forEach((el) => {
        const plainText = document.createTextNode(el.textContent);
        el.parentNode.replaceChild(plainText, el);
    });
}

document.addEventListener("keydown", function(event) {
    if (event.key === "Escape") {
        removeMarkedText();
    }
    const button = document.querySelector('.aa-DetachedSearchButton');
    // Check if CMD+K was pressed (on Mac)
    if (event.metaKey && event.keyCode === 75) {
        // Execute your function here
        console.log('CMD+K was pressed');
        button.click();
    }
    // Check if CTRL+K was pressed (on Windows/Linux)
    else if (event.ctrlKey && event.keyCode === 75) {
        // Execute your function here
        console.log('CTRL+K was pressed');
        button.click();
    }
});
