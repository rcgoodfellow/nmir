//nmir visualization environment

var camera, scene, renderer, theme;
var net;

function initCanvas() {

  var width = window.innerWidth;
  var height = window.innerHeight;
  theme = {
    background: 0xe0e0e0,
    node: 0x7ca8ef,
    endpoint: 0x1a3a6d,
    line: 0x142744
  }

  scene = new THREE.Scene();
  camera = new THREE.OrthographicCamera( 
    width/-2,
    width/2,
    height/2,
    height/-2,
    -500, 1000
  );

  renderer = new THREE.WebGLRenderer({
    antialias: true
  });
  renderer.setSize( window.innerWidth, window.innerHeight );
  renderer.setClearColor( theme.background, 1 );
  renderer.setPixelRatio( window.devicePixelRatio );
  renderer.setSize( width, height );
  document.body.appendChild( renderer.domElement );


  camera.position.z = 100;

}

function addNode(data, x, y) {

  var geometry = new THREE.CircleGeometry( 15, 32 );
  var material = new THREE.MeshBasicMaterial( { color: theme.node } );
  var node = new THREE.Mesh( geometry, material );
  node.data = data;
  node.name = data.id;
  node.position.x = x;
  node.position.y = y;
  scene.add( node );

  data.endpoints.forEach((endpoint, i, es) => {
    addEndpoint(node, endpoint);
  });

}

function addEndpoint(node, data) {

  var geometry = new THREE.CircleGeometry( 3, 16 );
  var material = new THREE.MeshBasicMaterial( { color: theme.endpoint } );
  var endpoint = new THREE.Mesh( geometry, material );
  endpoint.data = data;
  endpoint.name = data.id;
  node.add( endpoint );

}

function addLinkLine(data, x, y) {

  var material = new THREE.LineBasicMaterial({ 
    linewidth: 2,
    opacity: 0.9,
    color: theme.line 
  });
  var geometry = new THREE.Geometry();
  geometry.vertices.push(x.parent.position);
  geometry.vertices.push(y.parent.position);
  var line = new THREE.Line(geometry, material);
  scene.add( line );

}

var render = function () {
  requestAnimationFrame( render );
  camera.lookAt( scene.position );
  renderer.render(scene, camera);
};

function initData() {

  /*
  $.getJSON("https://mirror.deterlab.net/nmir/4net.json", (json) => {
    console.log(json);
    net = json;
    loadData();
  });
  */

  net = topo;
  loadData();


}

function loadData() {

  net.nodes.forEach((node, i, ns) => {
    //var r = 100;
    //var y = r*Math.sin(i*(Math.PI/6));
    //var x = r*Math.cos(i*(Math.PI/6));
    //console.log("adding node "+node.props.name);
    addNode(node, node.props.position.x, node.props.position.y);
  });

  net.links.forEach((link, i, ls) => {
    link.endpoints[0].forEach((a, i, es) => {
      link.endpoints[1].forEach((b, i, es) => {
        var na = scene.getObjectByName(a.id);
        var nb = scene.getObjectByName(b.id);
        addLinkLine(link, na, nb);
      });
    });
  });
  
  console.log(scene);

}

$(document).ready(() => {

  initCanvas();
  initData();
  render();

})

