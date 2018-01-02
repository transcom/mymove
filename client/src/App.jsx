import React, { Component } from 'react';
import Feedback from './scenes/Feedback/Feedback';
import Header from './shared/Header/Header';
import Footer from './shared/Footer/Footer';
import './App.css';

class App extends Component {
  render() {
    return (
      <div className="App site">
        <Header />
        <main className="site__content">
          <Feedback />
        </main>
        <Footer />
      </div>
    );
  }
}

export default App;
