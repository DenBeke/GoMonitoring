var sock = null;
var wsuri = "ws://localhost:8000/ws";



var lineChartData = {
    labels : ["","","","","","",""],
    datasets : [
        {
            fillColor : "rgba(151,187,205,0.5)",
            strokeColor : "rgba(151,187,205,1)",
            pointColor : "rgba(151,187,205,1)",
            pointStrokeColor : "#fff",
            data : [0,0,0,0,0,0,0]
        }
    ]
};

function initLine() {
    var options = {
        animation : false,
        scaleOverride : true,
        scaleSteps : 10,//Number - The number of steps in a hard coded scale
        scaleStepWidth : 10,//Number - The value jump in the hard coded scale				
        scaleStartValue : 0,//Number - The scale starting value
        responsive: true,
        maintainAspectRatio: true,
    };
    
    //var ctx = $("#canvas").get(0).getContext("2d");
    var myLine = new Chart(document.getElementById('canvas').getContext("2d")).Line( lineChartData, options );
}

/*
socket.on('pushdata', function (data) {
    self.lineChartData().datasets[0].data.shift();
    self.lineChartData().datasets[0].data.push(data);
    
    self.initLine();
});
*/



window.onload = function() {

    self.initLine();

    sock = new WebSocket(wsuri);

    sock.onopen = function() {
        console.log("connected to " + wsuri);
    }

    sock.onclose = function(e) {
        console.log("connection closed (" + e.code + ")");
    }

    sock.onmessage = function(e) {
        result = JSON.parse(e.data);
        console.log(result.UsedPercent);
        $('#cpu').html(result.UsedPercent + "%");
        
        lineChartData.datasets[0].data.shift();
        lineChartData.datasets[0].data.push(result.UsedPercent);
        
        self.initLine();
        
        
        //$('<li>' + e.data '</li>').hide().prependTo('#messages').fadeIn(1000);
    }
};