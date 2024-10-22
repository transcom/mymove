import React from 'react';
import { Route } from 'react-router-dom';
import { every, some, get, findKey, pick } from 'lodash';

import { generalRoutes, customerRoutes } from 'constants/routes';
import generatePath from 'shared/WizardPage/generatePath';
import { NULL_UUID } from 'shared/constants';
import BackupContact from 'pages/MyMove/Profile/BackupContact';
import ProfileReview from 'scenes/Review/ProfileReview';
import Home from 'pages/MyMove/Home';
import DodInfo from 'pages/MyMove/Profile/DodInfo';
import SMName from 'pages/MyMove/Profile/Name';
import ContactInfo from 'pages/MyMove/Profile/ContactInfo';
import Orders from 'pages/MyMove/Orders';
import UploadOrders from 'pages/MyMove/UploadOrders';
import SelectShipmentType from 'pages/MyMove/SelectShipmentType';
import BackupAddress from 'pages/MyMove/Profile/BackupAddress';
import ResidentialAddress from 'pages/MyMove/Profile/ResidentialAddress';
import Review from 'pages/MyMove/Review/Review';
import Agreement from 'pages/MyMove/Agreement';
import ValidationCode from 'pages/MyMove/Profile/ValidationCode';

const PageNotInFlow = () => (
  <div className="usa-grid">
    <h1>Missing Context</h1>
    You are trying to load a page that the system does not have context for. Please go to the home page and try again.
  </div>
);

const always = () => true;
const never = () => false;
const myFirstRodeo = (props) => !props.lastMoveIsCanceled;
const notMyFirstRodeo = (props) => props.lastMoveIsCanceled;
const inGhcFlow = (props) => props.context.flags.ghcFlow;
const isCurrentMoveSubmitted = ({ move }) => {
  return get(move, 'status', 'DRAFT') === 'SUBMITTED';
};

const pages = {
  [customerRoutes.VALIDATION_CODE_PATH]: {
    isInFlow: myFirstRodeo,
    isComplete: ({ sm }) => sm.is_profile_complete || every([sm.edipi, sm.affiliation]),
    render: () => <ValidationCode />,
  },
  [customerRoutes.DOD_INFO_PATH]: {
    isInFlow: myFirstRodeo,
    isComplete: ({ sm }) => sm.is_profile_complete || every([sm.edipi, sm.affiliation]),
    render: () => <DodInfo />,
  },
  [customerRoutes.NAME_PATH]: {
    isInFlow: myFirstRodeo,
    isComplete: ({ sm }) => sm.is_profile_complete || every([sm.first_name, sm.last_name]),
    render: () => <SMName />,
  },
  [customerRoutes.CONTACT_INFO_PATH]: {
    isInFlow: myFirstRodeo,
    isComplete: ({ sm }) =>
      sm.is_profile_complete ||
      (every([sm.telephone, sm.personal_email]) && some([sm.phone_is_preferred, sm.email_is_preferred])),
    render: () => <ContactInfo />,
  },
  [customerRoutes.CURRENT_ADDRESS_PATH]: {
    isInFlow: myFirstRodeo,
    isComplete: ({ sm }) => sm.is_profile_complete || Boolean(sm.residential_address),
    render: () => <ResidentialAddress />,
  },
  [customerRoutes.BACKUP_ADDRESS_PATH]: {
    isInFlow: myFirstRodeo,
    isComplete: ({ sm }) => sm.is_profile_complete || Boolean(sm.backup_mailing_address),
    render: () => <BackupAddress />,
  },
  [customerRoutes.BACKUP_CONTACTS_PATH]: {
    isInFlow: myFirstRodeo,
    isComplete: ({ sm, orders, move, ppm, backupContacts }) => {
      return sm.is_profile_complete || backupContacts.length > 0;
    },
    render: () => <BackupContact />,
    description: 'Backup contacts',
  },
  [generalRoutes.HOME_PATH]: {
    isInFlow: (props) => {
      return myFirstRodeo(props) && inGhcFlow(props);
    },
    isComplete: never,
    render: (key, pages) => {
      return <Home />;
    },
  },
  '/profile-review': {
    isInFlow: notMyFirstRodeo,
    isComplete: always,
    render: (key, pages) => <ProfileReview pages={pages} pageKey={key} />,
  },
  [customerRoutes.ORDERS_INFO_PATH]: {
    isInFlow: always,
    isComplete: ({ sm, orders }) =>
      every([
        orders.orders_type,
        orders.issue_date,
        orders.report_by_date,
        get(orders, 'new_duty_location.id', NULL_UUID) !== NULL_UUID,
        get(orders, 'origin_duty_location.id', NULL_UUID) !== NULL_UUID,
        orders.grade,
      ]),
    render: (key, pages) => <Orders />,
  },
  [customerRoutes.ORDERS_UPLOAD_PATH]: {
    isInFlow: always,
    isComplete: ({ sm, orders, uploads }) =>
      get(orders, 'uploaded_orders.uploads', []).length > 0 || uploads.length > 0,
    render: (key, pages, description, props) => <UploadOrders />,
    description: 'Upload your orders',
  },
  [customerRoutes.SHIPMENT_SELECT_TYPE_PATH]: {
    isInFlow: always,
    isComplete: ({ sm, orders, move }) => get(move, 'mtoShipments', []).length > 0,
    render: () => <SelectShipmentType />,
  },
  [customerRoutes.MOVE_REVIEW_PATH]: {
    isInFlow: always,
    isComplete: ({ sm, orders, move, ppm, mtoShipment }) => isCurrentMoveSubmitted(move),
    render: (props) => {
      return <Review {...props} />;
    },
  },
  [customerRoutes.MOVE_AGREEMENT_PATH]: {
    isInFlow: always,
    isComplete: ({ sm, orders, move, ppm, mtoShipment }) => isCurrentMoveSubmitted(move),
    render: () => <Agreement />,
  },
};

export const getPagesInFlow = ({ move, conusStatus, lastMoveIsCanceled, context }) =>
  Object.keys(pages).filter((pageKey) => {
    const page = pages[pageKey];
    return page.isInFlow({ move, conusStatus, lastMoveIsCanceled, context });
  });

export const getNextIncompletePage = ({
  conusStatus = '',
  lastMoveIsCanceled = false,
  serviceMember = {},
  orders = {},
  uploads = [],
  move = {},
  ppm = {},
  mtoShipment = {},
  backupContacts = [],
  context = {},
  excludeHomePage = false,
}) => {
  excludeHomePage && delete pages['/'];
  const rawPath = findKey(
    pages,
    (p) =>
      p.isInFlow({ move, conusStatus, lastMoveIsCanceled, context }) &&
      !p.isComplete({ sm: serviceMember, orders, uploads, move, ppm, mtoShipment, backupContacts }),
  );
  const compiledPath = generatePath(rawPath, {
    serviceMemberId: get(serviceMember, 'id'),
    moveId: get(move, 'id'),
    orderId: get(orders, 'id'),
  });
  return compiledPath;
};

export const getWorkflowRoutes = (props) => {
  const flowProps = pick(props, ['move', 'conusStatus', 'lastMoveIsCanceled', 'context']);
  const pageList = getPagesInFlow(flowProps);
  return Object.keys(pages).map((key) => {
    const currPage = pages[key];
    if (currPage.isInFlow(flowProps)) {
      const render = currPage.render(key, pageList, currPage.description, props);
      return <Route end path={key} key={key} element={render} />;
    }
    return <Route end path={key} key={key} element={<PageNotInFlow />} />;
  });
};
