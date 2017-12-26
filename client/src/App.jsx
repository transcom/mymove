import React, { Component } from 'react';
import Feedback from './scenes/Feedback/Feedback';
import Header from './shared/Header/Header';
import './App.css';

class App extends Component {
  render() {
    return (
      <div className="App">
        <Header />
        <Feedback />
      </div>
    );
  }
}

export default App;
