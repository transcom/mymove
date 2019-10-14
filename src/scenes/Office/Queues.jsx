import React, {Component, Suspense} from "react";
import {RetrieveMovesForOffice} from "./api";
import QueueList from "./QueueList";
import QueueTable from "./QueueTable"


export class Queues extends Component {
  render() {
    return (
      <div className="usa-grid grid-wide queue-columns">
        <div className="queue-menu-column">
          <Suspense fallback={<div>Loading...</div>}>
            <QueueList />
          </Suspense>
        </div>
        <div className="queue-list-column">
          <Suspense fallback={<div>Loading...</div>}>
            <QueueTable queueType={this.props.match.params.queueType} retrieveMoves={RetrieveMovesForOffice} />
          </Suspense>
        </div>
      </div>
    );
  }
}

export default Queues;