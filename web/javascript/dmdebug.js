$(document).ready(function(){

    html = "<div class='debug'>"+
        "<span class='title'>Time Spent</span>"+
        "<span>Total: </span><span>"+dmtime.total+"ms</span> "+
        "<span>Query: </span><span>"+dmtime.query+"ms</span> "+
        "<span>Template: </span><span>"+dmtime.template+"ms</span>"+
        "</div>";
    $("body").append( html );
});
