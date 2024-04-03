import React from 'react';
import PropTypes from 'prop-types';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useNavigate } from 'react-router';
import { GridContainer } from '@trussworks/react-uswds';

import CustomerContactInfoForm from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm';
import { CUSTOMER, ORDERS } from 'constants/queryKeys';
import { servicesCounselingRoutes } from 'constants/routes';
import { updateCustomerInfo } from 'services/ghcApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { CustomerShape } from 'types';

const CreateMoveCustomerInfo = ({ customer, isLoading, isError, ordersId, onUpdate }) => {
  const navigate = useNavigate();

  const handleBack = () => {
    navigate('/');
  };
  const handleClose = () => {
    navigate(`../${servicesCounselingRoutes.ORDERS_ADD_PATH}`);
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
      queryClient.invalidateQueries([ORDERS, ordersId]);
      onUpdate('success');
      handleClose();
    },
    onError: () => {
      onUpdate('error');
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const onSubmit = (values) => {
    const {
      firstName,
      lastName,
      customerTelephone,
      customerEmail,
      customerAddress,
      suffix,
      middleName,
      name,
      email,
      telephone,
      backupAddress,
      phoneIsPreferred,
      emailIsPreferred,
      secondaryPhone,
    } = values;

    const body = {
      first_name: firstName,
      last_name: lastName,
      phone: customerTelephone,
      email: customerEmail,
      current_address: customerAddress,
      suffix,
      middle_name: middleName,
      backup_contact: {
        name,
        email,
        phone: telephone,
      },
      backupAddress,
      phoneIsPreferred,
      emailIsPreferred,
      secondaryTelephone: secondaryPhone || null,
    };
    mutateCustomerInfo({ customerId: customer.id, ifMatchETag: customer.eTag, body });
  };
  const initialValues = {
    firstName: customer.first_name || '',
    lastName: customer.last_name || '',
    middleName: customer.middle_name || '',
    suffix: customer.suffix || '',
    customerTelephone: customer.phone || '',
    customerEmail: customer.email || '',
    name: customer.backup_contact.name || '',
    telephone: customer.backup_contact.phone || '',
    secondaryPhone: customer.secondaryTelephone || '',
    email: customer.backup_contact.email || '',
    customerAddress: customer.current_address || '',
    backupAddress: customer.backupAddress || '',
    emailIsPreferred: customer.emailIsPreferred || false,
    phoneIsPreferred: customer.phoneIsPreferred || false,
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

CreateMoveCustomerInfo.propTypes = {
  customer: CustomerShape.isRequired,
  isLoading: PropTypes.bool.isRequired,
  isError: PropTypes.bool.isRequired,
  ordersId: PropTypes.string.isRequired,
  onUpdate: PropTypes.func.isRequired,
};

export default CreateMoveCustomerInfo;
