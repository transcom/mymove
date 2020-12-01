import React, { Component, Suspense } from 'react';
import { RetrieveMovesForOffice } from './api';
import QueueList from './QueueList';
import QueueTable from './QueueTable';
import LoadingPlaceholder from '../../shared/LoadingPlaceholder';

export class Queues extends Component {
  render() {
    const queueType = this.props.queueType || this.props.match.params.queueType;

    return (
      <div className="usa-grid grid-wide queue-columns">
        <div className="queue-menu-column">
          <Suspense fallback={<LoadingPlaceholder />}>
            <QueueList />
          </Suspense>
        </div>
        <div className="queue-list-column">
          <Suspense fallback={<LoadingPlaceholder />}>
            <QueueTable queueType={queueType} retrieveMoves={RetrieveMovesForOffice} />
          </Suspense>
        </div>
      </div>
    );
  }
}

export default Queues;
