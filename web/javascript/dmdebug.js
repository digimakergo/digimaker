function load(){
    html = "<div class='debug'>";
    if( errorLog != "" )
    {
        html += "<div class='info hide'>"+errorLog+"</div>";
    }
    html += "<div class='time'><a href='javascript:switchLog();' class='title' title='More debug info'><i class='fas fa-chevron-up'></i> <span>Time Spent</span></a>"+
            "<span>Total: </span><span>"+dmtime.total+"ms</span> "+
            "<span>Query: </span><span>"+dmtime.query+"ms</span> "+
            "<span>Template: </span><span>"+dmtime.template+"ms</span></div>";

    html += "</div>";
    $("body").append( html );
};

$(document).ready(function(){
load();

});

function switchLog(){
    var info = $( '.debug .info' );
    var icon = $( ".debug .title i" );
    if( info.is(":visible") ){
        info.fadeOut();
        icon.removeClass("fa-chevron-down");
        icon.addClass("fa-chevron-up");

    }else {
        info.fadeIn();
        icon.removeClass("fa-chevron-up");
        icon.addClass("fa-chevron-down");
    }
}
