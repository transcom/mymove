import { Alert, Label, TextInput } from '@trussworks/react-uswds';
import React, { useState } from 'react';
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

import styles from './RejectedOfficeUserShow.module.scss';

const RejectedOfficeUserShowTitle = () => {
  const record = useRecordContext();

  return <span>{`${record?.firstName} ${record?.lastName}`}</span>;
};

const RejectedOfficeUserShowRoles = () => {
  const record = useRecordContext();
  if (!record?.roles) return <p>This user has not rejected any roles.</p>;

  return (
    <ArrayField source="roles">
      <span>Rejected roles:</span>
      <Datagrid bulkActionButtons={false}>
        <TextField source="roleName" />
      </Datagrid>
    </ArrayField>
  );
};

// renders server and rej reason alerts
// renders approve/reject/edit buttons
// handles logic of approving/rejecting user
const RejectedOfficeUserActionButtons = () => {
  const [serverError] = useState('');
  const [rejectionReason, setRejectionReason] = useState('');
  const [rejectionReasonCheck, setRejectionReasonCheck] = useState('');

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
        <Label>Rejection reason (required if rejecting)</Label>
        <TextInput
          label="Rejection reason"
          source="rejectionReason"
          value={rejectionReason}
          onChange={(e) => {
            setRejectionReason(e.target.value);
            // removing error banner if text is entered
            setRejectionReasonCheck('');
          }}
        />
      </div>
    </>
  );
};

const RejectedOfficeUserShow = () => {
  return (
    <Show title={<RejectedOfficeUserShowTitle />}>
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
        <RejectedOfficeUserShowRoles />
        <ReferenceField label="Transportation Office" source="transportationOfficeId" reference="offices" sortBy="name">
          <TextField component="pre" source="name" />
        </ReferenceField>
        <DateField label="Account rejected at" source="createdAt" showTime />
      </SimpleShowLayout>
      <RejectedOfficeUserActionButtons />
    </Show>
  );
};

export default RejectedOfficeUserShow;
