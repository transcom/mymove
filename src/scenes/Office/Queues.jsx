import React, { Suspense } from 'react';
import { useParams } from 'react-router-dom';
import { RetrieveMovesForOffice } from './api';
import QueueList from './QueueList';
import QueueTable from './QueueTable';
import LoadingPlaceholder from '../../shared/LoadingPlaceholder';

const Queues = (props) => {
  const params = useParams();
  const queueType = props.queueType || params.queueType;

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
};

export default Queues;
