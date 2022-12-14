import * as React from 'react';
import Map, {Marker} from 'react-map-gl';
import 'mapbox-gl/dist/mapbox-gl.css';

import scooter from './img/scooter.png';
import trotinette from './img/trotinette.png';
import bike from './img/bike.png';

import { w3cwebsocket as W3CWebSocket } from "websocket";

const client = new W3CWebSocket('ws://localhost:8080/ws');

function Icon(props) {
  if (props.type === 1) {
    return <Marker latitude={props.latitude} longitude={props.longitude}><img src={scooter} alt="scooter" width="20" height="20" /></Marker>;
  } else if (props.type === 2) {
    return <Marker latitude={props.latitude} longitude={props.longitude}><img src={trotinette} alt="trotinette" width="20" height="20" /></Marker>;
  } else if (props.type === 0) {
    return <Marker latitude={props.latitude} longitude={props.longitude}><img src={bike} alt="bike" width="20" height="20" /></Marker>;
  }
}

function MapIcons(props) {
  const icons = props.data;
  const listMarkers = Object.keys(icons).map((key) =>
    <Icon key={key} type={icons[key].type} longitude={icons[key].latitude} latitude={icons[key].longitude} />
  );

  console.log(listMarkers);

  return <div> {listMarkers} </div>
}

export default function App() {
  const [state, setState] = React.useState({});

  React.useEffect(() => {
    client.onopen = () => {
      console.log('WebSocket Client Connected');
    };
    client.onmessage = (message) => {
      const dataFromServer = JSON.parse(message.data);

      const tmp = {type: dataFromServer.VType, longitude: dataFromServer.Longitude, latitude: dataFromServer.Latitude};
      
      setState(prevState => ({
        ...prevState,
        [dataFromServer.Id]: tmp
      }));

    };
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
      <MapIcons data={state} />
    </Map>
  );
}

