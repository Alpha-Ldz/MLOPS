import * as React from 'react';
import Map, {Marker} from 'react-map-gl';
import 'mapbox-gl/dist/mapbox-gl.css';
import scooter from './img/scooter.png';
import trotinette from './img/trotinette.png';
import bike from './img/bike.png';


function App() {
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
    </Map>
  );
}

export default App;
