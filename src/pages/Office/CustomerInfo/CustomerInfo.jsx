/* eslint-disable camelcase */
import React from 'react';
import PropTypes from 'prop-types';
import { useHistory, useParams } from 'react-router-dom';
import { queryCache, useMutation } from 'react-query';
import { GridContainer } from '@trussworks/react-uswds';

import CustomerContactInfoForm from '../../../components/Office/CustomerContactInfoForm/CustomerContactInfoForm';

import { updateCustomerInfo } from 'services/ghcApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { CustomerShape } from 'types/order';

const CustomerInfo = ({ customer, isLoading, isError }) => {
  const { moveCode } = useParams();
  const history = useHistory();

  const handleClose = () => {
    history.push(`/counseling/moves/${moveCode}/details`);
  };

  const [mutateCustomerInfo] = useMutation(updateCustomerInfo, {
    onSuccess: (data, variables) => {
      // TODO: cache stuff
      // const updatedOrder = data.orders[variables.orderID];
      // queryCache.setQueryData([ORDERS, variables.orderID], {
      //   orders: {
      //     [`${variables.orderID}`]: updatedOrder,
      //   },
      // });
      // queryCache.invalidateQueries(ORDERS);
      handleClose();
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      // TODO: Handle error some how
      // RA Summary: eslint: no-console - System Information Leak: External
      // RA: The linter flags any use of console.
      // RA: This console displays an error message from unsuccessful mutation.
      // RA: TODO: As indicated, this error needs to be handled and needs further investigation and work.
      // RA: POAM story here: https://dp3.atlassian.net/browse/MB-5597
      // RA Developer Status: Known Issue
      // RA Validator Status: Known Issue
      // RA Modified Severity: CAT II
      // eslint-disable-next-line no-console
      console.log(errorMsg);
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const onSubmit = (values) => {
    const { firstName, lastName, customerTelephone, customerEmail, customerAddress, name, email, telephone } = values;
    const body = {
      first_name: firstName,
      last_name: lastName,
      phone: customerTelephone,
      email: customerEmail,
      current_address: customerAddress,
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
    middleName: customer.middle_name, // TODO: not sure this is implemented on backend
    suffix: customer.suffix, // TODO: not sure this is implemented on backend
    customerTelephone: customer.phone,
    customerEmail: customer.email,
    name: customer.backup_contact.name,
    telephone: customer.backup_contact.phone,
    email: customer.backup_contact.email,
    customerAddress: customer.current_address,
  };

  return (
    <GridContainer>
      <h1>Customer Info</h1>
      <CustomerContactInfoForm initialValues={initialValues} onBack={handleClose} onSubmit={onSubmit} />
    </GridContainer>
  );
};

CustomerInfo.propTypes = {
  customer: CustomerShape.isRequired,
  isLoading: PropTypes.bool.isRequired,
  isError: PropTypes.bool.isRequired,
};
export default CustomerInfo;
