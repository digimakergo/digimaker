$(document).ready(function(){

    html = "<div class='debug'>"+
        "<div><a href='#' class='title' title='More debug info'><i class='fas fa-chevron-up'></i> <span>Time Spent</span></a>"+
        "<span>Total: </span><span>"+dmtime.total+"ms</span> "+
        "<span>Query: </span><span>"+dmtime.query+"ms</span> "+
        "<span>Template: </span><span>"+dmtime.template+"ms</span></div>"+
        "</div>";
    $("body").append( html );
});
