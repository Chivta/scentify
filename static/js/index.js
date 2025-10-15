document.addEventListener("DOMContentLoaded", function() {
    const form = document.getElementById("generateForm");
    form.addEventListener("submit",function(event){
        event.preventDefault();
        handleFormSubmit();
    });
});


function handleFormSubmit(){
    const value = document.getElementById("request").value;
    
    fetch("/generate",{
        method: "POST",
        headers: { "Content-Type": "text" },
        body: value
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