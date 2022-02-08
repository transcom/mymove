import React from 'react';
import { string } from 'prop-types';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useGHCGetMoveHistory } from 'hooks/queries';

const MoveHistory = ({ moveCode }) => {
  const { moveHistory, isLoading, isError } = useGHCGetMoveHistory(moveCode);
  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;
  return (
    <div className="grid-container-desktop-lg" data-testid="move-history">
      <h1>Move History</h1>
      <div className="container">
        <div>
          <pre>{JSON.stringify(moveHistory, null, 2)}</pre>
        </div>
      </div>
    </div>
  );
};

MoveHistory.propTypes = {
  moveCode: string.isRequired,
};

export default MoveHistory;
