var page = require('webpage').create();
page.onConsoleMessage = function (msg) {
    console.log(msg);
};

page.open('http://oris.orientacnisporty.cz/Adresar', function () {
    page.render('page.png');
    page.evaluate(function () {

        var clubs = [];

        while (true) {
            var list = document.querySelectorAll("tr");
            for (var i = 2; i < list.length; i++) { // skip table header
                var name = list[i].children[0].children[0].text;
                clubs.push(name);
            }
            
            var next = document.querySelector("a.paginate_enabled_next");

            if (!!next) {
                next.click();
            } else {
                break;
            }
        }
        console.log(clubs.length);
        for (var i = 0; i < clubs.length; i++) {
            console.log(clubs[i]);
        }

    });
    phantom.exit();
});
