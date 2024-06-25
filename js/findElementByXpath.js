return (function findElementByXpath(xPath) {
  const xpath = `"${xPath[0]}"`;
  const el = document.evaluate(xpath, document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue;
  return el
}).apply(null, arguments)
