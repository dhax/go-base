let audioChunks = [];
let recorder;

const recordButton = document.getElementById("record");
const stopButton = document.getElementById("stop");
const audioElement = document.getElementById("audio");

recordButton.onclick = () => {
    navigator.mediaDevices.getUserMedia({ audio: true })
        .then(stream => {
            recorder = new MediaRecorder(stream);
            recorder.ondataavailable = e => {
                audioChunks.push(e.data);
                if (recorder.state == "inactive") {
                    const blob = new Blob(audioChunks, { type: "audio/mpeg" });
                    const url = URL.createObjectURL(blob);
                    audioElement.src = url;
                    uploadAudio(blob);
                }
            };
            recorder.start();
        })
        .catch(error => console.log(error));
}

stopButton.onclick = () => {
    recorder.stop();
}

function uploadAudio(blob) {
    const xhr = new XMLHttpRequest();
    xhr.open("POST", "/upload", true);
    xhr.setRequestHeader("Content-Type", "audio/mpeg");
    xhr.onreadystatechange = () => {
        if (xhr.readyState === 4 && xhr.status === 200) {
            console.log("Audio uploaded successfully");
        }
    };
    xhr.send(blob);
}
