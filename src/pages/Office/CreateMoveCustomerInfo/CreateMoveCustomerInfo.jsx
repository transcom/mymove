import React from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { generatePath, useNavigate, useParams } from 'react-router';
import { GridContainer } from '@trussworks/react-uswds';

import CustomerContactInfoForm, {
  backupContactName,
} from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm';
import { CUSTOMER } from 'constants/queryKeys';
import { servicesCounselingRoutes } from 'constants/routes';
import { updateCustomerInfo } from 'services/ghcApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useCustomerQuery } from 'hooks/queries';
import { milmoveLogger } from 'utils/milmoveLog';
import { formatTrueFalseInputValue } from 'utils/formatters';

const CreateMoveCustomerInfo = () => {
  const { customerId } = useParams();
  const { customerData, isLoading, isError } = useCustomerQuery(customerId);
  const navigate = useNavigate();

  const handleBack = () => {
    navigate('/');
  };
  const handleClose = () => {
    navigate(generatePath(servicesCounselingRoutes.BASE_CUSTOMERS_ORDERS_ADD_PATH, { customerId }));
  };
  const queryClient = useQueryClient();
  const { mutate: mutateCustomerInfo } = useMutation(updateCustomerInfo, {
    onSuccess: (data, variables) => {
      const updatedCustomer = data.customer[variables.customerId];
      queryClient.setQueryData([CUSTOMER, variables.customerId], {
        customer: {
          [`${variables.customerId}`]: updatedCustomer,
        },
      });
      queryClient.invalidateQueries([CUSTOMER, variables.customerId]);
      handleClose();
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLogger.error(errorMsg);
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const onSubmit = (values) => {
    const cacUser = values.cacUser === 'true';
    const {
      firstName,
      lastName,
      customerTelephone,
      customerEmail,
      customerAddress,
      suffix,
      middleName,
      backupAddress,
      phoneIsPreferred,
      emailIsPreferred,
      secondaryPhone,
    } = values;

    const backupFirstName = (values[backupContactName.toString()]?.firstName || '').trim();
    const backupLastName = (values[backupContactName.toString()]?.lastName || '').trim();
    const backupFullName = `${backupFirstName} ${backupLastName}`.trim();

    const body = {
      first_name: firstName,
      last_name: lastName,
      phone: customerTelephone,
      email: customerEmail,
      current_address: customerAddress,
      suffix,
      middle_name: middleName,
      backup_contact: {
        name: backupFullName,
        email: values[backupContactName.toString()]?.email || '',
        phone: values[backupContactName.toString()]?.telephone || '',
      },
      backupAddress,
      phoneIsPreferred,
      emailIsPreferred,
      secondaryTelephone: secondaryPhone || null,
      cac_validated: cacUser,
    };
    mutateCustomerInfo({ customerId: customerData.id, ifMatchETag: customerData.eTag, body });
  };

  const backupContactFullName = (customerData?.backup_contact?.name || '').trim();
  const [backupContactFirstName = '', backupContactLastName = ''] = backupContactFullName
    .split(/ (.+)/)
    .map((part) => part?.trim())
    .filter(Boolean);

  const initialValues = {
    firstName: customerData?.first_name || '',
    lastName: customerData?.last_name || '',
    middleName: customerData?.middle_name || '',
    suffix: customerData?.suffix || '',
    customerTelephone: customerData?.phone || '',
    customerEmail: customerData?.email || '',
    secondaryPhone: customerData?.secondaryTelephone || '',
    customerAddress: customerData?.current_address || '',
    backupAddress: customerData?.backupAddress || '',
    emailIsPreferred: customerData?.emailIsPreferred || false,
    phoneIsPreferred: customerData?.phoneIsPreferred || false,
    cacUser: formatTrueFalseInputValue(customerData?.cacValidated),
    [backupContactName]: {
      firstName: backupContactFirstName,
      lastName: backupContactLastName,
      email: customerData?.backup_contact?.email || '',
      telephone: customerData?.backup_contact?.phone || '',
    },
  };

  return (
    <div>
      <GridContainer>
        <h1>Confirm Customer Info</h1>
        <CustomerContactInfoForm initialValues={initialValues} onBack={handleBack} onSubmit={onSubmit} />
      </GridContainer>
    </div>
  );
};

export default CreateMoveCustomerInfo;
