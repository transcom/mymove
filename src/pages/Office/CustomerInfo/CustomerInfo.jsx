/* eslint-disable camelcase */
import React from 'react';
import PropTypes from 'prop-types';
import { queryCache, useMutation } from 'react-query';
import { useNavigate } from 'react-router-dom';
import { GridContainer } from '@trussworks/react-uswds';

import CustomerContactInfoForm from '../../../components/Office/CustomerContactInfoForm/CustomerContactInfoForm';

import styles from './CustomerInfo.module.scss';

import { CUSTOMER, ORDERS } from 'constants/queryKeys';
import { servicesCounselingRoutes } from 'constants/routes';
import { updateCustomerInfo } from 'services/ghcApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { CustomerShape } from 'types/order';

const CustomerInfo = ({ customer, isLoading, isError, ordersId, onUpdate }) => {
  const navigate = useNavigate();

  const handleClose = () => {
    navigate(`../${servicesCounselingRoutes.MOVE_VIEW_PATH}`);
  };

  const [mutateCustomerInfo] = useMutation(updateCustomerInfo, {
    onSuccess: (data, variables) => {
      const updatedCustomer = data.customer[variables.customerId];
      queryCache.setQueryData([CUSTOMER, variables.customerId], {
        customer: {
          [`${variables.customerId}`]: updatedCustomer,
        },
      });
      queryCache.invalidateQueries([CUSTOMER, variables.customerId]);
      queryCache.invalidateQueries([ORDERS, ordersId]);
      onUpdate('success');
      handleClose();
    },
    onError: () => {
      // TODO: Handle error some how - see https://dp3.atlassian.net/browse/MB-5597
      onUpdate('error');
      handleClose();
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
    };
    mutateCustomerInfo({ customerId: customer.id, ifMatchETag: customer.eTag, body });
  };

  const initialValues = {
    firstName: customer.first_name,
    lastName: customer.last_name,
    middleName: customer.middle_name,
    suffix: customer.suffix,
    customerTelephone: customer.phone,
    customerEmail: customer.email,
    name: customer.backup_contact.name,
    telephone: customer.backup_contact.phone,
    email: customer.backup_contact.email,
    customerAddress: customer.current_address,
  };

  return (
    <div className={styles.customerInfoPage}>
      <GridContainer>
        <h1>Customer Info</h1>
        <CustomerContactInfoForm initialValues={initialValues} onBack={handleClose} onSubmit={onSubmit} />
      </GridContainer>
    </div>
  );
};

CustomerInfo.propTypes = {
  customer: CustomerShape.isRequired,
  isLoading: PropTypes.bool.isRequired,
  isError: PropTypes.bool.isRequired,
  ordersId: PropTypes.string.isRequired,
  onUpdate: PropTypes.func.isRequired,
};
export default CustomerInfo;
