// window.alert("Hello")

let url = window.location.origin.replace("http","ws")
let socket = new WebSocket(url+"/ws")

console.log("請拍照後，沒問題再按保存");

socket.onopen = function(){
    console.log("連線到:" + socket);
}

socket.close = function(e){
    console.log("連線失敗 (" + "錯誤代碼:"+ e.code + ")");
}

socket.onmessage = function(e){
    console.log("Server回傳:"+e.data)
    var a = "data:image/jpeg;base64,"
    var b = a + e.data
    document.getElementById('img').src = b
    var img = document.getElementById('img');
    img.style.width = 400;


}

socket.onerror = function(e){
    console.log("Server錯誤:" + e.data);
}

function run(){
    var msg = document.getElementById('run').value;
    socket.send(msg);
};
function save(){
    var msg = document.getElementById('save').value;
    socket.send(msg);
    
};
