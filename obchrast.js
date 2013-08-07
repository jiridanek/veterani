var tables = document.querySelectorAll('[debug=band]');

for (var i=1; i < tables.length; i++) {
  var table = tables[i];
  var CATEGORY = table.previousElementSibling.firstElementChild.innerText.trim();
  var rows = table.querySelectorAll('[height="14"]');
  for (var j=0; j<rows.length; j++) {
    var row = rows[j];
    var cells = row.querySelectorAll('td');
    
    var jmeno = cells[1].innerText.trim();
    var id = cells[2].innerText.trim();
    var cas = cells[9].innerText.trim();

    if (cas === 'DISK') continue;
    
    var csv = [CATEGORY,jmeno,id,cas].join(',')
    console.log(csv)
  }
}