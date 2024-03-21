import { Alert, Button, Label, TextInput } from '@trussworks/react-uswds';
import React, { useState } from 'react';
import {
  ArrayField,
  Datagrid,
  DateField,
  EditButton,
  ReferenceField,
  Show,
  SimpleShowLayout,
  TextField,
  useRecordContext,
} from 'react-admin';
import { useNavigate } from 'react-router';

import styles from './RequestedOfficeUserShow.module.scss';

import { updateRequestedOfficeUser } from 'services/adminApi';
import { adminRoutes } from 'constants/routes';

const RequestedOfficeUserShowTitle = () => {
  const record = useRecordContext();

  return <span>{`${record?.firstName} ${record?.lastName}`}</span>;
};

const RequestedOfficeUserShowRoles = () => {
  const record = useRecordContext();
  if (!record?.roles) return <p>This user has not requested any roles.</p>;

  return (
    <ArrayField source="roles">
      <span>Requested roles:</span>
      <Datagrid bulkActionButtons={false}>
        <TextField source="roleName" />
      </Datagrid>
    </ArrayField>
  );
};

const RequestedOfficeUserActionButtons = () => {
  const [rejectionReason, setRejectionReason] = useState('');
  const [serverError, setServerError] = useState('');
  const navigate = useNavigate();
  const record = useRecordContext();

  const approve = async (user) => {
    const body = {
      edipi: user.edipi,
      firstName: user.firstName,
      middleInitials: user.middleInitials,
      lastName: user.lastName,
      otherUniqueId: user.otherUniqueId,
      rejectionReason: null,
      roles: user.roles,
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

  const reject = async (user, rejectionReasonInput) => {
    const body = {
      edipi: user.edipi,
      firstName: user.firstName,
      middleInitials: user.middleInitials,
      lastName: user.lastName,
      otherUniqueId: user.otherUniqueId,
      rejectionReason: rejectionReasonInput,
      roles: user.roles,
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
  };

  return (
    <>
      {serverError && (
        <Alert type="error" slim className={styles.error}>
          {serverError}
        </Alert>
      )}
      <div className={styles.rejectionInput}>
        <Label>Rejection reason (optional)</Label>
        <TextInput
          label="Rejection reason"
          source="rejectionReason"
          value={rejectionReason}
          onChange={(e) => setRejectionReason(e.target.value)}
        />
      </div>
      <div className={styles.btnContainer}>
        <Button
          className={styles.approveBtn}
          onClick={async () => {
            await approve(record);
          }}
        >
          Approve
        </Button>
        <Button
          className={styles.rejectBtn}
          onClick={async () => {
            await reject(record);
          }}
        >
          Reject
        </Button>
        <EditButton />
      </div>
    </>
  );
};

const RequestedOfficeUserShow = () => {
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
        <RequestedOfficeUserShowRoles />
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
