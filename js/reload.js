function reload() {
  // var container = document.getElementById("bodyWrapper");
  var container = document.getElementsByTagName("<placeholder>")[0];
  var content = container.innerHTML;
  container.innerHTML = content; 
  
  //this line is to watch the result in console , you can remove it later	
  console.log("Refreshed"); 
}

reload();