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
    const description = document.getElementById("request").value;
    const noteAmount = parseInt(document.getElementById("note-amount-slider").value);
    const silliness = parseInt(document.getElementById("silliness-level-slider").value);

    const body = {
        description: description,
        silliness: silliness,
        noteAmount: noteAmount,
    }

    fetch("/generate",{
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body)
    })
    .then(response => response.json())
    .then(data => {
        let cards = document.getElementById("cards");
        cards.innerHTML="";
        data.forEach(element => {
            cards.insertAdjacentHTML("beforeend",`
                <div class="card">
                    <div class="scentImage">  
                        <img src="${element.image}">
                    </div>
                    <small>${element.note}</small>
                    <div class="remove">
                        <img src="static/svg/x-mark.svg" class="xmark" alt="Remove">
                    </div>
                </div>`
            )
        });
    })
}

