document.addEventListener("DOMContentLoaded", function() {
    document.getElementById("generateForm")
    .addEventListener("submit",function(event){
        event.preventDefault();
        handleFormSubmit();
    });
    document.getElementById("note-amount-slider")
    .addEventListener("input", sliderTextUpdater("note-amount-slider","selected-note-amount"));
    document.getElementById("silliness-level-slider")
    .addEventListener("input", sliderTextUpdater("silliness-level-slider","selected-silliness-level"));

    let colorPicker = document.getElementById("generator-color-picker");
    colorPicker.addEventListener("input",()=>{
        document.getElementById("generator").style.backgroundColor=colorPicker.value;
    });

    document.getElementById("color-reset-button")
    .addEventListener("click",()=>{
        colorPicker.value="#cef8cf";
        document.getElementById("generator").style.backgroundColor="#cef8cf";
    });

    document.getElementById("dark-mode-switcher")
    .addEventListener("click", () => {
        const body = document.body;
        body.classList.toggle("dark-mode-body");
    });
});

function sliderTextUpdater(sliderId, textId) {
    return () => 
    document.getElementById(textId).textContent =
    document.getElementById(sliderId).value;
}

function handleFormSubmit(){
    let submitButton = document.getElementById("submit-button");
    if (submitButton.disabled) return;

    submitButton.disabled = true;
    const prevText = submitButton.value;
    submitButton.value = "Generating...";

    const description = document.getElementById("request").value;
    const noteAmount = parseInt(document.getElementById("note-amount-slider").value);
    const silliness = parseInt(document.getElementById("silliness-level-slider").value);
    const generateImages = document.getElementById("generate-images-checkbox").checked;

    let cards = document.getElementById("cards");
    cards.innerHTML="";

    const body = {
        description: description,
        silliness: silliness,
        noteAmount: noteAmount,
        generateImages: generateImages
    }

    fetch("/generate",{
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body)
    })
    .then(response => response.json())
    .then(data => {
        data.forEach(element => {
            cards.insertAdjacentHTML("beforeend",`
                <div class="card">
                    <div class="scentImage">  
                        <img src="${element.image}">
                    </div>
                    <small>${element.note}</small>
                    <div class="remove">
                        <img src="static/x-mark.svg" class="xmark">
                    </div>
                </div>`
            )
        });
    }).catch(err => {
        console.error(err);
    })
    .finally(() => {
        submitButton.disabled = false;
        submitButton.value = prevText;
    });
}

