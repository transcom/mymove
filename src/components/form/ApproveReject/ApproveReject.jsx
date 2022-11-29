import React from 'react';
import { func, bool, string } from 'prop-types';
import { Button, Fieldset, FormGroup, Label, Radio, Textarea } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './ApproveReject.module.scss';

const ApproveReject = ({
  id,
  approvedStatus,
  deniedStatus,
  currentStatus,
  handleApprovalChange,
  handleRejectCancel,
  handleRejectChange,
  handleChange,
  handleFormReset,
  rejectionReason,
  requestComplete,
  canEditRejection,
  setCanEditRejection,
  rejectOptionText,
}) => {
  return (
    <Fieldset className={styles.ApproveReject}>
      <div className={classnames(styles.statusOption, { [styles.selected]: currentStatus === approvedStatus })}>
        <Radio
          id={`approve-${id}`}
          checked={currentStatus === approvedStatus}
          value={approvedStatus}
          name="status"
          label="Approve"
          onChange={handleApprovalChange}
          data-testid="approveRadio"
        />
      </div>
      <div className={classnames(styles.statusOption, { [styles.selected]: currentStatus === deniedStatus })}>
        <Radio
          id={`reject-${id}`}
          checked={currentStatus === deniedStatus}
          value={deniedStatus}
          name="status"
          label={rejectOptionText}
          onChange={handleChange}
          data-testid="rejectRadio"
        />

        {currentStatus === deniedStatus && (
          <FormGroup>
            <Label htmlFor={`rejectReason-${id}`}>Reason for rejection</Label>
            {!canEditRejection && (
              <>
                <p data-testid="rejectionReasonReadOnly">{rejectionReason}</p>
                <Button
                  type="button"
                  unstyled
                  data-testid="editReasonButton"
                  className={styles.clearStatus}
                  onClick={() => setCanEditRejection(true)}
                  aria-label="Edit reason button"
                >
                  <span className="icon">
                    <FontAwesomeIcon icon="pen" title="Edit reason" alt="" />
                  </span>
                  <span aria-hidden="true">Edit reason</span>
                </Button>
              </>
            )}

            {!requestComplete && canEditRejection && (
              <>
                <Textarea
                  id={`rejectReason-${id}`}
                  name="rejectionReason"
                  onChange={handleChange}
                  value={rejectionReason}
                />
                <div className={styles.rejectionButtonGroup}>
                  <Button
                    id="rejectionSaveButton"
                    type="button"
                    data-testid="rejectionSaveButton"
                    onClick={handleRejectChange}
                    disabled={!rejectionReason}
                    aria-label="Rejection save button"
                  >
                    Save
                  </Button>
                  <Button
                    data-testid="cancelRejectionButton"
                    secondary
                    onClick={handleRejectCancel}
                    type="button"
                    aria-label="Cancel rejection button"
                  >
                    Cancel
                  </Button>
                </div>
              </>
            )}
          </FormGroup>
        )}
      </div>

      {(currentStatus === approvedStatus || currentStatus === deniedStatus) && (
        <Button
          type="button"
          unstyled
          data-testid="clearStatusButton"
          className={styles.clearStatus}
          onClick={handleFormReset}
          aria-label="Clear status"
        >
          <span className="icon">
            <FontAwesomeIcon icon="times" title="Clear status" alt=" " />
          </span>
          <span aria-hidden="true">Clear selection</span>
        </Button>
      )}
    </Fieldset>
  );
};

ApproveReject.propTypes = {
  id: string.isRequired,
  currentStatus: string,
  rejectionReason: string,
  rejectOptionText: string,
  requestComplete: bool,
  canEditRejection: bool,
  approvedStatus: string,
  deniedStatus: string,
  handleApprovalChange: func.isRequired,
  handleRejectCancel: func.isRequired,
  handleRejectChange: func.isRequired,
  handleChange: func.isRequired,
  handleFormReset: func.isRequired,
  setCanEditRejection: func.isRequired,
};

ApproveReject.defaultProps = {
  currentStatus: '',
  approvedStatus: '',
  deniedStatus: '',
  rejectionReason: '',
  rejectOptionText: 'Reject',
  requestComplete: false,
  canEditRejection: false,
};

export default ApproveReject;
