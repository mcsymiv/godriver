return (function setValueToActive(args) {
  document.activeElement.value=`"${args[0]}"`;
}).apply(null, arguments)
