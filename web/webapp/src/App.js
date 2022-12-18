import * as React from 'react';
import Map, {Marker} from 'react-map-gl';
import 'mapbox-gl/dist/mapbox-gl.css';

import scooter from './img/scooter.png';
import trotinette from './img/trotinette.png';
import bike from './img/bike.png';

import { w3cwebsocket as W3CWebSocket } from "websocket";

const client = new W3CWebSocket('ws://localhost:8080/ws');

function Icon(props) {

  props = JSON.parse(props.props);


  if (props.VType === 1) {
    return <Marker latitude={props.Longitude} longitude={props.Latitude}><img src={scooter} alt="scooter" width="20" height="20" /></Marker>;
  } else if (props.VType === 2) {
    return <Marker latitude={props.Longitude} longitude={props.Latitude}><img src={trotinette} alt="trotinette" width="20" height="20" /></Marker>;
  } else if (props.VType === 0) {
    return <Marker latitude={props.Longitude} longitude={props.Latitude}><img src={bike} alt="bike" width="20" height="20" /></Marker>;
  }
}

function MapIcons(props) {
  const icons = props.data;


  const listMarkers = Object.keys(icons).map((key) =>
    <Icon key={key} props={icons[key]}/>
  );

  return <div> {listMarkers} </div>
}

export default function App() {
  const [message, setMessage] = React.useState('');

  React.useEffect(() => {
    const interval = setInterval(() => {
      fetch('http://localhost:3333/random')
        .then(response => response.text())
        .then(message => setMessage(JSON.parse(message)));
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  return (
    <Map
      initialViewState={{
        latitude: 48.856416650506375,
        longitude: 2.3494125832340313,
        zoom: 11
      }}
      style={{width: "100vw", height: "100vh"}}
      mapStyle="mapbox://styles/mapbox/streets-v9"
      mapboxAccessToken={process.env.REACT_APP_MAPBOX_ACCESS_TOKEN}
    >
      <MapIcons data={message} />
    </Map>
  );
}

