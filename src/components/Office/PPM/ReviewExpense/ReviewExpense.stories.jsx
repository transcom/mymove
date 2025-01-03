import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import PPMShipmentInfo from '../ppmTestData';

import ReviewExpense from './ReviewExpense';

import { expenseTypes } from 'constants/ppmExpenseTypes';
import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components / PPM / Review Expense',
  component: ReviewExpense,
  decorators: [
    (Story) => (
      <MockProviders>
        <GridContainer>
          <Grid row>
            <Grid col desktop={{ col: 2, offset: 8 }}>
              <Story />
            </Grid>
          </Grid>
        </GridContainer>
      </MockProviders>
    ),
  ],
  argTypes: { onClose: { action: 'back button clicked' } },
};

const Template = (args) => <ReviewExpense {...args} />;

const documentSetsProps = [
  {
    documentSetType: 'WEIGHT_TICKET',
    documentSet: {
      adjustedNetWeight: null,
      allowableWeight: 3000,
      createdAt: '2024-05-16T18:50:33.689Z',
      eTag: 'MjAyNC0wNS0yMFQxNzo1MjowMS4xNTA3MDFa',
      emptyDocument: {
        id: '67739c19-37ca-4de2-8412-513b3564b70c',
        service_member_id: '3f13faa6-0cec-4842-887b-8a7a2d54797e',
        uploads: [
          {
            bytes: 5291,
            contentType: 'image/png',
            createdAt: '2024-05-16T18:50:45.453Z',
            filename: 'thumbnail_image001.png',
            id: '14bcfa31-8879-4ceb-83b5-d36f17005774',
            status: 'PROCESSING',
            updatedAt: '2024-05-16T18:50:45.453Z',
            url: '/storage/user/cebc022f-e133-4c6c-ad5c-e5db3c06c5e2/uploads/14bcfa31-8879-4ceb-83b5-d36f17005774?contentType=image%2Fpng',
          },
        ],
      },
      emptyDocumentId: '67739c19-37ca-4de2-8412-513b3564b70c',
      emptyWeight: 1000,
      fullDocument: {
        id: '219675be-c530-4ea4-85ae-8eee92215849',
        service_member_id: '3f13faa6-0cec-4842-887b-8a7a2d54797e',
        uploads: [
          {
            bytes: 5291,
            contentType: 'image/png',
            createdAt: '2024-05-16T18:50:51.465Z',
            filename: 'thumbnail_image001.png',
            id: 'c1fb5f37-b633-46b4-a827-fffa8f7f4b92',
            status: 'PROCESSING',
            updatedAt: '2024-05-16T18:50:51.465Z',
            url: '/storage/user/cebc022f-e133-4c6c-ad5c-e5db3c06c5e2/uploads/c1fb5f37-b633-46b4-a827-fffa8f7f4b92?contentType=image%2Fpng',
          },
        ],
      },
      fullDocumentId: '219675be-c530-4ea4-85ae-8eee92215849',
      fullWeight: 4000,
      id: 'c6014838-ff2c-4803-ab3d-8cea23689c10',
      missingEmptyWeightTicket: false,
      missingFullWeightTicket: false,
      netWeightRemarks: null,
      ownsTrailer: true,
      ppmShipmentId: '09b1d087-f8e5-4f46-b75a-297cff7a4d33',
      proofOfTrailerOwnershipDocument: {
        id: 'c087764e-6095-43f0-afda-49cd486820f7',
        service_member_id: '3f13faa6-0cec-4842-887b-8a7a2d54797e',
        uploads: [
          {
            bytes: 5291,
            contentType: 'image/png',
            createdAt: '2024-05-16T18:51:01.436Z',
            filename: 'thumbnail_image001.png',
            id: 'b850d264-29fb-4723-84f5-4a6488e7f358',
            status: 'PROCESSING',
            updatedAt: '2024-05-16T18:51:01.436Z',
            url: '/storage/user/cebc022f-e133-4c6c-ad5c-e5db3c06c5e2/uploads/b850d264-29fb-4723-84f5-4a6488e7f358?contentType=image%2Fpng',
          },
        ],
      },
      proofOfTrailerOwnershipDocumentId: 'c087764e-6095-43f0-afda-49cd486820f7',
      reason: null,
      status: 'APPROVED',
      trailerMeetsCriteria: true,
      updatedAt: '2024-05-20T17:52:01.150Z',
      vehicleDescription: 'test',
    },
    uploads: [
      {
        bytes: 5291,
        contentType: 'image/png',
        createdAt: '2024-05-16T18:50:45.453Z',
        filename: 'thumbnail_image001.png',
        id: '14bcfa31-8879-4ceb-83b5-d36f17005774',
        status: 'PROCESSING',
        updatedAt: '2024-05-16T18:50:45.453Z',
        url: '/storage/user/cebc022f-e133-4c6c-ad5c-e5db3c06c5e2/uploads/14bcfa31-8879-4ceb-83b5-d36f17005774?contentType=image%2Fpng',
      },
      {
        bytes: 5291,
        contentType: 'image/png',
        createdAt: '2024-05-16T18:50:51.465Z',
        filename: 'thumbnail_image001.png',
        id: 'c1fb5f37-b633-46b4-a827-fffa8f7f4b92',
        status: 'PROCESSING',
        updatedAt: '2024-05-16T18:50:51.465Z',
        url: '/storage/user/cebc022f-e133-4c6c-ad5c-e5db3c06c5e2/uploads/c1fb5f37-b633-46b4-a827-fffa8f7f4b92?contentType=image%2Fpng',
      },
      {
        bytes: 5291,
        contentType: 'image/png',
        createdAt: '2024-05-16T18:51:01.436Z',
        filename: 'thumbnail_image001.png',
        id: 'b850d264-29fb-4723-84f5-4a6488e7f358',
        status: 'PROCESSING',
        updatedAt: '2024-05-16T18:51:01.436Z',
        url: '/storage/user/cebc022f-e133-4c6c-ad5c-e5db3c06c5e2/uploads/b850d264-29fb-4723-84f5-4a6488e7f358?contentType=image%2Fpng',
      },
    ],
    tripNumber: 0,
  },
];

const documentSetIndex = 0;

export const Blank = Template.bind({});
Blank.args = {
  ppmShipmentInfo: PPMShipmentInfo,
  documentSets: documentSetsProps,
  documentSetIndex,
  tripNumber: 1,
  ppmNumber: '1',
};

export const NonStorage = Template.bind({});
NonStorage.args = {
  ppmShipmentInfo: PPMShipmentInfo,
  documentSets: documentSetsProps,
  documentSetIndex,
  tripNumber: 1,
  ppmNumber: '1',
  categoryIndex: 1,
  expense: {
    movingExpenseType: expenseTypes.PACKING_MATERIALS,
    description: 'boxes, tape, bubble wrap',
    amount: 12345,
  },
};

export const Storage = Template.bind({});
Storage.args = {
  documentSets: documentSetsProps,
  documentSetIndex,
  ppmShipmentInfo: PPMShipmentInfo,
  tripNumber: 1,
  ppmNumber: '1',
  categoryIndex: 1,
  expense: {
    movingExpenseType: expenseTypes.STORAGE,
    description: 'Pack n store',
    amount: 12345,
    sitStartDate: '2022-12-15',
    sitEndDate: '2022-12-25',
    weightStored: 2000,
    sitLocation: 'ORIGIN',
  },
};
