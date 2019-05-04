$(document).ready(function(){
    html = "<div class='debug'>";
    if( errorLog != "" )
    {
        html += "<div class='error'>"+errorLog+"</div>";
    }
    html += "<div class='time'><a href='#' class='title' title='More debug info'><i class='fas fa-chevron-up'></i> <span>Time Spent</span></a>"+
            "<span>Total: </span><span>"+dmtime.total+"ms</span> "+
            "<span>Query: </span><span>"+dmtime.query+"ms</span> "+
            "<span>Template: </span><span>"+dmtime.template+"ms</span></div>";

    html += "</div>";
    $("body").append( html );
});