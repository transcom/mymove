import React, { useState, useEffect } from 'react';
import { Alert, Button, Label, TextInput } from '@trussworks/react-uswds';
import {
  ArrayField,
  Datagrid,
  DateField,
  ReferenceField,
  Show,
  SimpleShowLayout,
  TextField,
  useRecordContext,
} from 'react-admin';
import { useNavigate } from 'react-router';

import RequestedOfficeUserPrivilegeConfirm from './RequestedOfficeUserPrivilegeConfirm';
import styles from './RequestedOfficeUserShow.module.scss';

import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { updateRequestedOfficeUser } from 'services/adminApi';
import { adminRoutes } from 'constants/routes';
import { FEATURE_FLAG_KEYS } from 'shared/constants';
import { elevatedPrivilegeTypes } from 'constants/userPrivileges';

// Helper to filter only SUPERVISOR privileges from a privileges array
export function getFilteredPrivileges(privileges) {
  return (privileges || []).filter((priv) => priv.privilegeType === elevatedPrivilegeTypes.SUPERVISOR);
}

const RequestedOfficeUserShowTitle = () => {
  const record = useRecordContext();

  return <span>{`${record?.firstName} ${record?.lastName}`}</span>;
};

/**
 * Displays a list of requested roles or privileges for an office user.
 * If displaying privileges, only shows those with privilegeType 'SUPERVISOR'.
 * If no items are present after filtering, displays a message indicating none were requested.
 * Used in the Requested Office User Show view for both roles and privileges.
 */
const RequestedOfficeUserShowRolesPrivileges = ({ recordSource, recordLabel, recordField }) => {
  const record = useRecordContext();
  const sourceLabel = typeof recordSource === 'string' ? recordSource.toLowerCase() : '';
  if (!record?.[recordSource]) return <p>{`This user has not requested any ${sourceLabel}.`}</p>;
  let items = record[recordSource] || [];
  if (recordSource === 'privileges') {
    items = getFilteredPrivileges(items);
  }
  if (!items.length) return <p>{`This user has not requested any ${sourceLabel}.`}</p>;

  return (
    <ArrayField source={recordSource} record={{ ...record, [recordSource]: items }}>
      <span id={`${recordSource}-label`}>
        <strong>{recordLabel}:</strong>
      </span>
      <Datagrid bulkActionButtons={false} aria-labelledby={`${recordSource}-label`}>
        <TextField source={recordField} />
      </Datagrid>
    </ArrayField>
  );
};

// renders server and rej reason alerts
// renders approve/reject/edit buttons
// handles logic of approving/rejecting user
const RequestedOfficeUserActionButtons = () => {
  const [serverError, setServerError] = useState('');
  const [rejectionReason, setRejectionReason] = useState('');
  const [rejectionReasonCheck, setRejectionReasonCheck] = useState('');
  const [approveDialogOpen, setApproveDialogOpen] = useState(false);
  const [checkedPrivileges, setCheckedPrivileges] = useState([]);
  const navigate = useNavigate();
  const record = useRecordContext();
  const [isRequestAccountPrivilegesFF, setRequestAccountPrivilegesFF] = useState(false);

  useEffect(() => {
    isBooleanFlagEnabled(FEATURE_FLAG_KEYS.REQUEST_ACCOUNT_PRIVILEGES).then(setRequestAccountPrivilegesFF);
  }, []);
  // if approved here, all values are good, but we want to change status to APPROVED
  const approve = async (user) => {
    setRejectionReasonCheck('');
    const body = {
      email: user.email,
      edipi: user.edipi,
      firstName: user.firstName,
      middleInitials: user.middleInitials,
      lastName: user.lastName,
      otherUniqueId: user.otherUniqueId,
      rejectionReason: null,
      roles: user.roles,
      privileges: user.privileges,
      status: 'APPROVED',
      telephone: user.telephone,
      transportationOfficeId: user.transportationOfficeId,
    };
    updateRequestedOfficeUser(record.id, body)
      .then(() => {
        navigate(adminRoutes.HOME_PATH);
      })
      .catch((error) => {
        setServerError(error);
      });
  };

  // if rejected here, all values are good, but we want to change status to REJECTED
  const reject = async (user, rejectionReasonInput) => {
    if (!rejectionReasonInput || rejectionReasonInput === '') {
      setRejectionReasonCheck('Please provide a rejection reason.');
    } else {
      const body = {
        email: user.email,
        edipi: user.edipi,
        firstName: user.firstName,
        middleInitials: user.middleInitials,
        lastName: user.lastName,
        otherUniqueId: user.otherUniqueId,
        rejectionReason: rejectionReasonInput,
        roles: user.roles,
        privileges: user.privileges,
        status: 'REJECTED',
        telephone: user.telephone,
        transportationOfficeId: user.transportationOfficeId,
      };
      updateRequestedOfficeUser(record.id, body)
        .then(() => {
          navigate(adminRoutes.HOME_PATH);
        })
        .catch((error) => {
          setServerError(error);
        });
    }
  };

  // Handler for privilege confirmation dialog
  const handlePrivilegeConfirm = async () => {
    setApproveDialogOpen(false);
    // Only include checked privileges in approval, and only SUPERVISOR privilege
    const filteredPrivileges = getFilteredPrivileges(record.privileges);
    const approvedPrivileges = filteredPrivileges.filter((priv) => checkedPrivileges.includes(priv.id)) || [];
    await approve({ ...record, privileges: approvedPrivileges });
  };

  // Handler for Approve button click
  const handleOnClickApprove = () => {
    const filteredPrivileges = getFilteredPrivileges(record.privileges);
    if (isRequestAccountPrivilegesFF && filteredPrivileges.length) {
      setCheckedPrivileges(filteredPrivileges.map((priv) => priv.id));
      setApproveDialogOpen(true);
      return;
    }
    approve(record);
  };

  return (
    <>
      {serverError && (
        <Alert type="error" slim className={styles.error}>
          {serverError}
        </Alert>
      )}
      {rejectionReasonCheck && (
        <Alert type="error" slim className={styles.error}>
          {rejectionReasonCheck}
        </Alert>
      )}
      <div className={styles.rejectionInput}>
        <Label htmlFor="show-rejection-reason-input">Rejection reason (required if rejecting)</Label>
        <TextInput
          id="show-rejection-reason-input"
          label="Rejection reason"
          aria-label="Rejection reason"
          source="rejectionReason"
          value={rejectionReason}
          onChange={(e) => {
            setRejectionReason(e.target.value);
            // removing error banner if text is entered
            setRejectionReasonCheck('');
          }}
        />
      </div>
      <div className={styles.btnContainer}>
        <Button
          className={styles.rejectBtn}
          onClick={async () => {
            await reject(record, rejectionReason);
          }}
        >
          Reject
        </Button>
        <Button className={styles.approveBtn} onClick={handleOnClickApprove}>
          Approve
        </Button>
      </div>
      <RequestedOfficeUserPrivilegeConfirm
        dialogId="show-approve-privilege-dialog"
        isOpen={approveDialogOpen}
        privileges={record?.privileges || []}
        checkedPrivileges={checkedPrivileges}
        setCheckedPrivileges={setCheckedPrivileges}
        onConfirm={handlePrivilegeConfirm}
        onClose={() => setApproveDialogOpen(false)}
      />
    </>
  );
};

const RequestedOfficeUserShow = () => {
  const [isRequestAccountPrivilegesFF, setRequestAccountPrivilegesFF] = useState(false);
  useEffect(() => {
    isBooleanFlagEnabled(FEATURE_FLAG_KEYS.REQUEST_ACCOUNT_PRIVILEGES).then(setRequestAccountPrivilegesFF);
  }, []);

  return (
    <Show title={<RequestedOfficeUserShowTitle />}>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="userId" label="User Id" />
        <TextField source="status" />
        <TextField source="email" />
        <TextField source="firstName" />
        <TextField source="middleInitials" />
        <TextField source="lastName" />
        <TextField source="telephone" />
        <TextField source="edipi" label="DODID#" />
        <TextField source="otherUniqueId" label="Other unique Id" />
        <RequestedOfficeUserShowRolesPrivileges
          recordSource="roles"
          recordLabel="Requested roles"
          recordField="roleName"
        />
        {isRequestAccountPrivilegesFF && (
          <RequestedOfficeUserShowRolesPrivileges
            recordSource="privileges"
            recordLabel="Requested privileges"
            recordField="privilegeName"
          />
        )}
        <ReferenceField label="Transportation Office" source="transportationOfficeId" reference="offices" sortBy="name">
          <TextField component="pre" source="name" />
        </ReferenceField>
        <DateField label="Account requested at" source="createdAt" showTime />
      </SimpleShowLayout>
      <RequestedOfficeUserActionButtons />
    </Show>
  );
};

export default RequestedOfficeUserShow;
