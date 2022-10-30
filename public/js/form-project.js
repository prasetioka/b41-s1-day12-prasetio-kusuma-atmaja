let dataProject = [];

function addProject(event) {
    event.preventDefault();

    let title = document.getElementById("input-title").value;


    let start = document.getElementById("start-date").value;
    let end = document.getElementById("end-date").value;
    const dateOne = new Date(start);
    const dateTwo = new Date(end);
    const time = Math.floor(dateTwo-dateOne);
    const days = Math.floor(time / (1000*60*60*24));
    const months = Math.floor(days / 30);

    let description = document.getElementById("project-description").value;
    let technologies = document.getElementById("tech-list").value;
    let image = document.getElementById("upload-image").files[0];
    image = URL.createObjectURL(image);

    let project = {
        title,
        months,
        description,
        technologies,
        image,
        postAt: new Date()
    };

    dataProject.push(project);

    console.log(start);
    console.log(end);
    console.log(dataProject);

    renderProject();
    
    setInterval (function() {
        renderProject()
    }, 1000)
}

function renderProject() {
    document.getElementById("contents").innerHTML = '';
    
    for (let index = 0; index < dataProject.length; index++) {
        // console.log("test", dataProject[index]);

        document.getElementById("contents").innerHTML += `
            <div class="card-project">
                <a href="../pages/my-project-detail.html">
                <img src="${dataProject[index].image}">
                </a>
                <h2>${dataProject[index].title}</h2>
                <h3>durasi: ${dataProject[index].months} bulan</h3>
                <p>${dataProject[index].description}</p>
                <div class="card-icon">
                    <ul>
                        <li><a href="#"><i class="fa-brands fa-google-play"></i></a></li>
                        <li><a href="#"><i class="fa-brands fa-android"></i></a></li>
                        <li><a href="#"><i class="fa-brands fa-java"></i></a></li>
                    </ul>
                </div>
                <div class="card-button">
                    <button>edit</button>
                    <button>delete</button>
                </div>
                <p class="timer">${getDistanceTime(dataProject[index].postAt)}</p>
            </div>`
    }
}

function getDistanceTime(time) {
    let timeNow = new Date()
    let timePost = time
    
    let distance = timeNow - timePost
    
    let milisecond = 1000
    let secondInHours = 3600
    let hoursInDay = 24
    
    let distanceDay = Math.floor(distance / (milisecond * secondInHours * hoursInDay))
    
    let distanceHours = Math.floor(distance / (milisecond * 60 * 60))
    
    let distanceMinutes = Math.floor(distance / (milisecond * 60))
    
    let distanceSecond = Math.floor(distance / milisecond)
    
    if (distanceDay > 0) {
        return `${distanceDay} day ago`
    } else if (distanceHours > 0) {
        return `${distanceHours} hour(s) ago`
    } else if (distanceMinutes > 0) {
        return `${distanceMinutes} minute(s) ago`
    } else {
        return `${distanceSecond} second(s) ago`
    }
}