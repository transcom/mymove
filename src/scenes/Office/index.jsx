import React, { Component } from 'react';

class QueueList extends Component {
  render() {
    return <div style={{ background: 'rgb(255,200,200)' }}>QueueList</div>;
  }
}

class QueueTable extends Component {
  render() {
    return <div style={{ background: 'rgb(200,255,200)' }}>QueueTable</div>;
  }
}

class QueueHeader extends Component {
  render() {
    return <div style={{ background: 'rgb(200,200,255)' }}>QueueHeader</div>;
  }
}

class OfficeWrapper extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Office';
  }

  render() {
    return (
      <div className="Office site">
        <main className="site__content">
          <div>
            <div class="usa-grid">
              <QueueHeader />
            </div>
            <div class="usa-grid">
              <div class="usa-width-one-fourth">
                <QueueList />
              </div>
              <div class="usa-width-three-fourths">
                <QueueTable />
              </div>
            </div>
          </div>
        </main>
      </div>
    );
  }
}

export default OfficeWrapper;
