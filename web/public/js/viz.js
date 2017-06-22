//nmir visualization environment

var camera, scene, renderer, theme;
var net;

function initCanvas() {

  var width = window.innerWidth;
  var height = window.innerHeight;
  theme = {
    background: 0xe0e0e0,

    node_color: 0x7ca8ef,
    node_size: 10,

    endpoint_color: 0x1a3a6d,

    line_color: 0x142744,
    line_opacity: 0.9,
    line_width: 2,

    label_offset: 10
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

function addNode(data, x, y, g) {

  var node = new THREE.Group();
  var geometry = new THREE.CircleGeometry( theme.node_size, 32 );
  var material = new THREE.MeshBasicMaterial( { color: theme.node_color } );
  var body = new THREE.Mesh( geometry, material );
  body.position.z = 5;
  node.add(body);

  node.data = data;
  node.name = data.id;
  node.position.x = x;
  node.position.y = y;
  g.add( node );

  data.endpoints.forEach((endpoint, i, es) => {
    addEndpoint(node, endpoint);
  });
  
  addLabel(node)

}

function htmlXY(node) {

  var vector = node.position.clone();
  var canvas = renderer.domElement;

  // map to normalized device coordinate (NDC) space
  node.parent.updateMatrixWorld();
  vector.setFromMatrixPosition(node.matrixWorld);
  vector.project( camera );
  
  // map to 2D screen space
  vector.x = Math.round( (   vector.x + 1 ) * canvas.width  / 2 );
  vector.y = Math.round( ( - vector.y + 1 ) * canvas.height / 2 );
  vector.z = 0;

  return vector

}

function addLabel(node) {

  var vec = htmlXY(node);
  x = vec.x;
  y = vec.y;

  theta = node.data.props.label_angle;
  x += theme.label_offset*Math.cos(theta)
  y -= theme.label_offset*Math.sin(theta)

  var text = node.data.props.name;

  var label = document.createElement('div');
  label.classList.add('nodelabel');
  label.style.position = 'absolute';
  label.style.zIndex = 5;
  label.innerHTML = text;
  label.style.top = y + 'px';
  label.style.left = x + 'px';
  document.body.appendChild(label);

  h = label.clientHeight;
  w = label.clientWidth;

  label.style.top = (y - h) + 'px';
}

function addEndpoint(node, data) {

  var geometry = new THREE.CircleGeometry( 3, 16 );
  var material = new THREE.MeshBasicMaterial( { color: theme.endpoint_color } );
  var endpoint = new THREE.Mesh( geometry, material );
  endpoint.data = data;
  endpoint.name = data.id;
  endpoint.z = 0;
  node.add( endpoint );

}

function addLinkLine(lnk, x, y, g) {

  var material = new THREE.LineBasicMaterial({ 
    linewidth: theme.line_width,
    opacity: theme.line_opacity,
    color: theme.line_color
  });
  var geometry = new THREE.Geometry();
  g.updateMatrixWorld();

  var xp = x.parent.position.clone();
  var yp = y.parent.position.clone();
  
  if(!lnk.props.local) {
    xp.setFromMatrixPosition(x.matrixWorld);
    yp.setFromMatrixPosition(y.matrixWorld);
  }

  geometry.vertices.push(xp);
  geometry.vertices.push(yp);
  var line = new THREE.Line(geometry, material);
  g.add( line );

}

var render = function () {
  requestAnimationFrame( render );
  camera.lookAt( scene.position );
  renderer.render(scene, camera);
};

function initData() {

  net = topo;
  loadData();

}

function loadData() {

  showNet(net, scene);
  
  console.log(scene);

}

function showNet(net, parent) {

  var g = new THREE.Group();
  g.position.x = net.props.position.x
  g.position.y = net.props.position.y
  parent.add( g );

  if(net.nets != null) {
    net.nets.forEach((n, i, ns) => {
      showNet(n, g)
    });
  }


  if(net.nodes != null) {
    net.nodes.forEach((node, i, ns) => {
      addNode(node, node.props.position.x, node.props.position.y, g);
    });
  }


  if(net.links != null) {
    net.links.forEach((link, i, ls) => {
      link.endpoints[0].forEach((a, i, es) => {
        link.endpoints[1].forEach((b, i, es) => {
          var na = parent.getObjectByName(a.id);
          var nb = parent.getObjectByName(b.id);
          addLinkLine(link, na, nb, g);
        });
      });
    });
  }


}

$(document).ready(() => {

  initCanvas();
  initData();
  render();

})

