import React from 'react';
import { Button } from '@trussworks/react-uswds';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';
import { useBulkAssignmentQueries } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

export const BulkAssignmentModal = ({ onSubmit, onClose, queueType }) => {
  const { bulkAssignmentData, isLoading, isError } = useBulkAssignmentQueries(queueType);
  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;
  return (
    <Modal>
      <ModalClose handleClick={() => onClose()} />
      <ModalTitle>
        <h3>Bulk Assignment</h3>
      </ModalTitle>
      <p>Available moves: {bulkAssignmentData?.availableMoves}</p>
      {bulkAssignmentData?.availableOfficeUsers?.map((user) => {
        return (
          <div>
            <span>{`${user.lastName}, ${user.firstName}   ||| workload: ${user.workload || 0}`}</span>
          </div>
        );
      })}
      <ModalActions autofocus="true">
        <Button
          data-focus="true"
          className="usa-button--destructive"
          type="submit"
          data-testid="modalSubmitButton"
          onClick={() => onSubmit()}
        >
          submit text
        </Button>
        <Button className="usa-button--secondary" type="button" onClick={() => onClose()} data-testid="modalBackButton">
          close text
        </Button>
      </ModalActions>
    </Modal>
  );
};

export default connectModal(BulkAssignmentModal);
