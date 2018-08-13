var oldOnload = window.onload;

window.onload = function() {
    if (typeof oldOnload == 'function') {
        oldOnload();
    }

    webSkill();
}

window.onresize = function() {
    if(window.innerWidth <= 600){
        let skillcon1 = document.getElementById("skills-web-con");
        skillcon1.style.display = "inline-flex";
        
        let skillcon2 = document.getElementById("skills-mobile-con");
        skillcon2.style.display = "inline-flex";
    }
    else {
        webSkill();
    }
}

function sendEmail() {
    let validation = 0;

    let email = (document.getElementById('email').value != "") ? 
                    document.getElementById('email').value : 
                    validation++;
    let subject = (document.getElementById('subject').value != "") ? 
                    document.getElementById('subject').value : 
                    validation++;
                    
    let fname = (document.getElementById('fname').value != "") ? 
                    document.getElementById('fname').value : 
                    validation++;
    
    let lname = (document.getElementById('lname').value != "") ? 
                    document.getElementById('lname').value : 
                    validation++;
                    
    let message = (document.getElementById('message').value != "") ? 
                    document.getElementById('message').value : 
                    validation++;

    var template_params = {
        "email": email,
        "reply_to": email,
        "subject": subject,
        "first_name": fname,
        "last_name": lname,
        "message": message
    }

    var data = {
        service_id: 'gmail',
        template_id: 'default_template_v1',
        user_id: 'user_Zy7dboUT6wWHWVHpMoHaI',
        template_params: template_params
    };

    let headers = {
        "Content-type": "application/json"
    };

    let options = {
        method: 'POST',
        headers: headers,
        body: JSON.stringify(data)
    };

    if (validation == 0){
        fetch('https://api.emailjs.com/api/v1.0/email/send', options)
        .then((httpResponse) => {
            if (httpResponse.ok) {
                sentMessage(true);
            } else {
                return httpResponse.text()
                    .then(text => Promise.reject(text));
            }
        })
        .catch((error) => {
            console.log(error);
            sentMessage(false);
        });
    }
    else {
        sentMessage(false);
    }
}

function skills() {
    document.getElementById("html5").style.width = "80%";
    document.getElementById("css3").style.width = "60%";
    document.getElementById("javascript").style.width = "50%";
    document.getElementById("android-java").style.width = "80%";
    document.getElementById("react-native").style.width = "60%";
}

function sentMessage(status){
    (status) ? 
        document.getElementById('message-success').style.display = "inline-flex" :
        document.getElementById('message-fail').style.display = "inline-flex";
}

function openNav() {
    document.getElementById("mySidenav").style.width = "100%";
}

function closeNav() {
    document.getElementById("mySidenav").style.width = "0";
}

function webSkill() {
    if(window.innerWidth >= 600){
        let oldskillcon = document.getElementById("skills-mobile-con");
        oldskillcon.style.display = "none";
        
        let skillcon = document.getElementById("skills-web-con");
        skillcon.style.display = "inline-flex";
    
        let oldskillbtn = document.getElementById("skills-mobile-btn");
        oldskillbtn.classList.remove("skill-select");
    
        let skillbtn = document.getElementById("skills-web-btn");
        skillbtn.classList.add("skill-select");
    }
}

function mobileSkill() {
    if(window.innerWidth >= 600){
        let oldskillcon = document.getElementById("skills-web-con");
        oldskillcon.style.display = "none";
        
        let skillcon = document.getElementById("skills-mobile-con");
        skillcon.style.display = "inline-flex";
    
        let oldskillbtn = document.getElementById("skills-web-btn");
        oldskillbtn.classList.remove("skill-select");
    
        let skillbtn = document.getElementById("skills-mobile-btn");
        skillbtn.classList.add("skill-select");
    }
}