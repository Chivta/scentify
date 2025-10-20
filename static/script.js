document.addEventListener("DOMContentLoaded", function() {
    loadLocalStorage();

    document.getElementById("generateForm")
    .addEventListener("submit",function(event){
        event.preventDefault();
        handleFormSubmit();
    });

    document.getElementById("note-amount-slider")
    .addEventListener("input", sliderTextUpdater("note-amount-slider","selected-note-amount"));
    
    document.getElementById("silliness-level-slider")
    .addEventListener("input", sliderTextUpdater("silliness-level-slider","selected-silliness-level"));


    let generateImagesCheckbox = document.getElementById("generate-images-checkbox");
    generateImagesCheckbox.addEventListener("change",()=>{
        localStorage.setItem("generate-images", generateImagesCheckbox.checked ? "true" : "false");
    });

    let colorPicker = document.getElementById("generator-color-picker");
    colorPicker.addEventListener("input",()=>{
        document.getElementById("generator").style.backgroundColor=colorPicker.value;
        localStorage.setItem("generator-color", colorPicker.value);
    });

    document.getElementById("color-reset-button")
    .addEventListener("click",()=>{
        colorPicker.value="#cef8cf";
        document.getElementById("generator").style.backgroundColor="#cef8cf";
        localStorage.setItem("generator-color", "#cef8cf");
    });

    document.getElementById("dark-mode-switcher")
    .addEventListener("click", toggleDarkMode);
});

function toggleDarkMode(){
    const isDark = document.body.classList.toggle("dark-mode-body");
    localStorage.setItem("dark-mode", isDark ? "true" : "false");
}

function loadLocalStorage(){
    if (localStorage.getItem("dark-mode") === "true"){
        document.body.classList.add("dark-mode-body");
    } else {
        document.body.classList.remove("dark-mode-body");
    }

    const genImgs = localStorage.getItem("generate-images");
    const genImgsCheckbox = document.getElementById("generate-images-checkbox");
    if (genImgsCheckbox) genImgsCheckbox.checked = (genImgs === "true");

    const noteVal = localStorage.getItem("selected-note-amount");
    const noteSlider = document.getElementById("note-amount-slider");
    const noteText = document.getElementById("selected-note-amount");
    if (noteSlider && noteVal !== null) noteSlider.value = noteVal;
    if (noteText && noteSlider) noteText.textContent = noteSlider.value;

    const sillVal = localStorage.getItem("selected-silliness-level");
    const sillSlider = document.getElementById("silliness-level-slider");
    const sillText = document.getElementById("selected-silliness-level");
    if (sillSlider && sillVal !== null) sillSlider.value = sillVal;
    if (sillText && sillSlider) sillText.textContent = sillSlider.value;

    const color = localStorage.getItem("generator-color");
    const colorPicker = document.getElementById("generator-color-picker");
    const gen = document.getElementById("generator");
    if (colorPicker && color !== null) colorPicker.value = color;
    if (gen && color !== null) gen.style.backgroundColor = color;
}

function sliderTextUpdater(sliderId, textId) {
    return () => {
        document.getElementById(textId).textContent =
        document.getElementById(sliderId).value;
        localStorage.setItem(textId,document.getElementById(sliderId).value);
    } 
    
}

function handleFormSubmit(){
    let submitButton = document.getElementById("submit-button");
    if (submitButton.disabled) return;
    const description = document.getElementById("request").value;
    if (description == "") return
    submitButton.disabled = true;
    const prevText = submitButton.value;
    submitButton.value = "Generating...";

    let loader = document.getElementById("loader");
    loader.hidden = false;

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
        loader.hidden = true;
    });
}

