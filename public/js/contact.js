function showData() {
    let showName = document.getElementById('input-name').value;
    let showEmail = document.getElementById('input-email').value;
    let showPhone = document.getElementById('input-phone').value;
    let showSubject = document.getElementById('input-subject').value;
    let showMessage = document.getElementById('input-message').value;

    console.log(showName);
    console.log(showEmail);
    console.log(showPhone);
  
    if (showName == '') {
        return alert("Name can't be empty.")
    } else if (showEmail == '') {
        return alert("Email can't be empty.")
    } else if (showPhone == '') {
        return alert("Phone number can't be empty.")
    } else if (showSubject == '') {
        return alert("Subject can't be empty.")
    } else if (showMessage == '') {
        return alert("Please leave a message.")
    }

    let emailReceiver = 'prasetiokusumaatmaja@gmail.com';

    let a = document.createElement('a'); // membuat tag <a></a>

    // a.href = `mailto:${emailReceiver}?subject:${showSubject}&body= Hello, My name is ${showName}, ${showMessage}`;

    a.href = `https://mail.google.com/mail/?view=cm&fs=1&to=${emailReceiver}&su=${showSubject}&body=Hello, My name is ${showName}.%0D%0A${showMessage}`;

    a.target = "_blank";

    a.click()
}