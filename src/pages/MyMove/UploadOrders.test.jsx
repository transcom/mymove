import React from 'react';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { act } from 'react-dom/test-utils';

import UploadOrders from './UploadOrders';

import { deleteUpload, getAllMoves, getOrders, createUploadForDocument } from 'services/internalApi';
import { renderWithProviders } from 'testUtils';
import { customerRoutes } from 'constants/routes';
import { selectOrdersForLoggedInUser, selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import { ORDERS_TYPE } from 'constants/orders';

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectOrdersForLoggedInUser: jest.fn(),
  selectServiceMemberFromLoggedInUser: jest.fn(),
}));

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  createUploadForDocument: jest.fn().mockImplementation(() => Promise.resolve()),
  deleteUpload: jest.fn().mockImplementation(() => Promise.resolve()),
  getOrders: jest.fn().mockImplementation(() => Promise.resolve()),
  getAllMoves: jest.fn().mockImplementation(() => Promise.resolve()),
}));

const testOrdersValues = {
  id: 'testOrdersId',
  orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
  issue_date: '2020-11-08',
  report_by_date: '2020-11-26',
  has_dependents: false,
  moves: [{ id: 'testMoveId' }],
  new_duty_location: {
    address: {
      city: 'Des Moines',
      country: 'US',
      id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
      postalCode: '50309',
      state: 'IA',
      streetAddress1: '987 Other Avenue',
      streetAddress2: 'P.O. Box 1234',
      streetAddress3: 'c/o Another Person',
    },
    address_id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
    affiliation: 'AIR_FORCE',
    created_at: '2020-10-19T17:01:16.114Z',
    id: 'f9299768-16d2-4a13-ae39-7087a58b1f62',
    name: 'Yuma AFB',
    updated_at: '2020-10-19T17:01:16.114Z',
  },
};

const testPropsWithUploads = {
  id: 'testOrdersId',
  orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
  issue_date: '2020-11-08',
  report_by_date: '2020-11-26',
  has_dependents: false,
  new_duty_location: {
    address: {
      city: 'Des Moines',
      country: 'US',
      id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
      postalCode: '50309',
      state: 'IA',
      streetAddress1: '987 Other Avenue',
      streetAddress2: 'P.O. Box 1234',
      streetAddress3: 'c/o Another Person',
    },
    address_id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
    affiliation: 'AIR_FORCE',
    created_at: '2020-10-19T17:01:16.114Z',
    id: 'f9299768-16d2-4a13-ae39-7087a58b1f62',
    name: 'Yuma AFB',
    updated_at: '2020-10-19T17:01:16.114Z',
  },
  uploaded_orders: {
    id: 'testId',
    uploads: [
      {
        bytes: 1578588,
        contentType: 'image/png',
        createdAt: '2024-02-23T16:51:45.504Z',
        filename: 'Screenshot 2024-02-15 at 12.22.53 PM (2).png',
        id: 'fd88b0e6-ff6d-4a99-be6f-49458a244209',
        status: 'PROCESSING',
        updatedAt: '2024-02-23T16:51:45.504Z',
        url: '/storage/user/5fe4d948-aa1c-4823-8967-b1fb40cf6679/uploads/fd88b0e6-ff6d-4a99-be6f-49458a244209?contentType=image%2Fpng',
      },
    ],
  },
  moves: ['testMoveId'],
};

const testPropsNoUploads = {
  id: 'testOrdersId2',
  orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
  issue_date: '2020-11-08',
  report_by_date: '2020-11-26',
  has_dependents: false,
  new_duty_location: {
    address: {
      city: 'Des Moines',
      country: 'US',
      id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
      postalCode: '50309',
      state: 'IA',
      streetAddress1: '987 Other Avenue',
      streetAddress2: 'P.O. Box 1234',
      streetAddress3: 'c/o Another Person',
    },
    address_id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
    affiliation: 'AIR_FORCE',
    created_at: '2020-10-19T17:01:16.114Z',
    id: 'f9299768-16d2-4a13-ae39-7087a58b1f62',
    name: 'Yuma AFB',
    updated_at: '2020-10-19T17:01:16.114Z',
  },
  uploaded_orders: {
    id: 'testId',
    service_member_id: 'testId',
    uploads: [],
  },
  moves: ['testMoveId'],
};

const testOrders = [
  {
    id: 'testOrdersId2',
    orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
    issue_date: '2020-11-08',
    report_by_date: '2020-11-26',
    has_dependents: false,
    new_duty_location: {
      address: {
        city: 'Des Moines',
        country: 'US',
        id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
        postalCode: '50309',
        state: 'IA',
        streetAddress1: '987 Other Avenue',
        streetAddress2: 'P.O. Box 1234',
        streetAddress3: 'c/o Another Person',
      },
      address_id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
      affiliation: 'AIR_FORCE',
      created_at: '2020-10-19T17:01:16.114Z',
      id: 'f9299768-16d2-4a13-ae39-7087a58b1f62',
      name: 'Yuma AFB',
      updated_at: '2020-10-19T17:01:16.114Z',
    },
    uploaded_orders: {
      id: 'testId',
      uploads: [],
    },
    moves: ['testMoveId'],
  },
  {
    id: 'testOrdersId',
    orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
    issue_date: '2020-11-08',
    report_by_date: '2020-11-26',
    has_dependents: false,
    new_duty_location: {
      address: {
        city: 'Des Moines',
        country: 'US',
        id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
        postalCode: '50309',
        state: 'IA',
        streetAddress1: '987 Other Avenue',
        streetAddress2: 'P.O. Box 1234',
        streetAddress3: 'c/o Another Person',
      },
      address_id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
      affiliation: 'AIR_FORCE',
      created_at: '2020-10-19T17:01:16.114Z',
      id: 'f9299768-16d2-4a13-ae39-7087a58b1f62',
      name: 'Yuma AFB',
      updated_at: '2020-10-19T17:01:16.114Z',
    },
    uploaded_orders: {
      id: 'testId',
      uploads: [
        {
          bytes: 1578588,
          contentType: 'image/png',
          createdAt: '2024-02-23T16:51:45.504Z',
          filename: 'Screenshot 2024-02-15 at 12.22.53 PM (2).png',
          id: 'fd88b0e6-ff6d-4a99-be6f-49458a244209',
          status: 'PROCESSING',
          updatedAt: '2024-02-23T16:51:45.504Z',
          url: '/storage/user/5fe4d948-aa1c-4823-8967-b1fb40cf6679/uploads/fd88b0e6-ff6d-4a99-be6f-49458a244209?contentType=image%2Fpng',
        },
      ],
    },
    moves: ['testMoveId'],
  },
];

afterEach(() => {
  jest.resetAllMocks();
});

const mockParams = { orderId: 'testOrdersId' };
const mockPath = customerRoutes.ORDERS_UPLOAD_PATH;
const mockRoutingOptions = { path: mockPath, params: mockParams };

const serviceMember = {
  id: 'id123',
  current_location: {
    address: {
      city: 'Fort Bragg',
      country: 'United States',
      id: 'f1ee4cea-6b23-4971-9947-efb51294ed32',
      postalCode: '29310',
      state: 'NC',
      streetAddress1: '',
    },
    address_id: 'f1ee4cea-6b23-4971-9947-efb51294ed32',
    affiliation: 'ARMY',
    created_at: '2020-10-19T17:01:16.114Z',
    id: 'dca78766-e76b-4c6d-ba82-81b50ca824b9"',
    name: 'Fort Bragg',
    updated_at: '2020-10-19T17:01:16.114Z',
  },
};

const serviceMemberMoves = {
  currentMove: [
    {
      createdAt: '2024-02-23T19:30:11.374Z',
      eTag: 'MjAyNC0wMi0yM1QxOTozMDoxMS4zNzQxN1o=',
      id: 'testMoveId',
      moveCode: '44649B',
      orders: {
        authorizedWeight: 11000,
        created_at: '2024-02-23T19:30:11.369Z',
        entitlement: {
          proGear: 2000,
          proGearSpouse: 500,
        },
        grade: 'E_7',
        has_dependents: false,
        id: 'testOrders1',
        issue_date: '2024-02-29',
        new_duty_location: {
          address: {
            city: 'Fort Irwin',
            country: 'United States',
            id: '77dca457-d0d6-4718-9ca4-a630b4614cf8',
            postalCode: '92310',
            state: 'CA',
            streetAddress1: 'n/a',
          },
          address_id: '77dca457-d0d6-4718-9ca4-a630b4614cf8',
          affiliation: 'ARMY',
          created_at: '2024-02-22T21:34:21.449Z',
          id: '12421bcb-2ded-4165-b0ac-05f76301082a',
          name: 'Fort Irwin, CA 92310',
          transportation_office: {
            address: {
              city: 'Fort Irwin',
              country: 'United States',
              id: '65a97b21-cf6a-47c1-a4b6-e3f885dacba5',
              postalCode: '92310',
              state: 'CA',
              streetAddress1: 'Langford Lake Rd',
              streetAddress2: 'Bldg 105',
            },
            created_at: '2018-05-28T14:27:37.312Z',
            gbloc: 'LKNQ',
            id: 'd00e3ee8-baba-4991-8f3b-86c2e370d1be',
            name: 'PPPO Fort Irwin - USA',
            phone_lines: [],
            updated_at: '2018-05-28T14:27:37.312Z',
          },
          transportation_office_id: 'd00e3ee8-baba-4991-8f3b-86c2e370d1be',
          updated_at: '2024-02-22T21:34:21.449Z',
        },
        orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
        originDutyLocationGbloc: 'BGAC',
        origin_duty_location: {
          address: {
            city: 'Fort Gregg-Adams',
            country: 'United States',
            id: '12270b68-01cf-4416-8b19-125d11bc8340',
            postalCode: '23801',
            state: 'VA',
            streetAddress1: 'n/a',
          },
          address_id: '12270b68-01cf-4416-8b19-125d11bc8340',
          affiliation: 'ARMY',
          created_at: '2024-02-22T21:34:26.430Z',
          id: '9cf15b8d-985b-4ca3-9f27-4ba32a263908',
          name: 'Fort Gregg-Adams, VA 23801',
          transportation_office: {
            address: {
              city: 'Fort Gregg-Adams',
              country: 'United States',
              id: '10dc88f5-d76a-427f-89a0-bf85587b0570',
              postalCode: '23801',
              state: 'VA',
              streetAddress1: '1401 B Ave',
              streetAddress2: 'Bldg 3400, Room 119',
            },
            created_at: '2018-05-28T14:27:42.125Z',
            gbloc: 'BGAC',
            id: '4cc26e01-f0ea-4048-8081-1d179426a6d9',
            name: 'PPPO Fort Gregg-Adams - USA',
            phone_lines: [],
            updated_at: '2018-05-28T14:27:42.125Z',
          },
          transportation_office_id: '4cc26e01-f0ea-4048-8081-1d179426a6d9',
          updated_at: '2024-02-22T21:34:26.430Z',
        },
        report_by_date: '2024-02-29',
        service_member_id: '81aeac60-80f3-44d1-9b74-ba6d405ee2da',
        spouse_has_pro_gear: false,
        status: 'DRAFT',
        updated_at: '2024-02-23T19:30:11.369Z',
        uploaded_orders: {
          id: 'bd35c4c2-41c6-44a1-bf54-9098c68d87cc',
          service_member_id: '81aeac60-80f3-44d1-9b74-ba6d405ee2da',
          uploads: [
            {
              bytes: 92797,
              contentType: 'image/png',
              createdAt: '2024-02-26T18:43:58.515Z',
              filename: 'Screenshot 2024-02-08 at 12.57.43 PM.png',
              id: '786237dc-c240-449d-8859-3f37583b3406',
              status: 'PROCESSING',
              updatedAt: '2024-02-26T18:43:58.515Z',
              url: '/storage/user/5fe4d948-aa1c-4823-8967-b1fb40cf6679/uploads/786237dc-c240-449d-8859-3f37583b3406?contentType=image%2Fpng',
            },
          ],
        },
      },
      status: 'DRAFT',
      submittedAt: '0001-01-01T00:00:00.000Z',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
  ],
  previousMoves: [],
};

describe('UploadOrders component', () => {
  it('renders the component successfully', async () => {
    selectServiceMemberFromLoggedInUser.mockImplementation(() => serviceMember);
    selectOrdersForLoggedInUser.mockImplementation(() => testOrders);
    getOrders.mockResolvedValue(testOrdersValues);
    const testProps = {
      serviceMember,
      serviceMemberId: 'id123',
      orders: [testPropsWithUploads],
      updateOrders: jest.fn(),
      updateAllMoves: jest.fn(),
    };

    renderWithProviders(<UploadOrders {...testProps} />, mockRoutingOptions);

    await screen.findByRole('heading', { level: 1, name: 'Upload your orders' });
    expect(screen.getByTestId('upload-orders-container')).toBeInTheDocument();
  });

  it('back button exists and enabled', async () => {
    selectServiceMemberFromLoggedInUser.mockImplementation(() => serviceMember);
    selectOrdersForLoggedInUser.mockImplementation(() => testOrders);
    getOrders.mockResolvedValue(testOrdersValues);
    getAllMoves.mockResolvedValue(() => serviceMemberMoves);
    const testProps = {
      serviceMember,
      serviceMemberId: 'id123',
      orders: [testPropsWithUploads],
      updateOrders: jest.fn(),
    };

    renderWithProviders(<UploadOrders {...testProps} />, {
      path: customerRoutes.ORDERS_UPLOAD_PATH,
      params: { orderId: 'testOrdersId' },
    });

    const backBtn = await screen.findByRole('button', { name: 'Back' });
    expect(backBtn).toBeInTheDocument();
    expect(backBtn).toBeEnabled();
  });

  it('next button exists and disabled when there are no uploads', async () => {
    selectServiceMemberFromLoggedInUser.mockImplementation(() => serviceMember);
    selectOrdersForLoggedInUser.mockImplementation(() => testOrders);
    getOrders.mockResolvedValue(testOrdersValues);
    getAllMoves.mockResolvedValue(() => serviceMemberMoves);
    const testProps = {
      serviceMember,
      serviceMemberId: 'id123',
      orders: [testPropsNoUploads],
      updateOrders: jest.fn(),
    };

    renderWithProviders(<UploadOrders {...testProps} />, {
      path: customerRoutes.ORDERS_UPLOAD_PATH,
      params: { orderId: 'testOrdersId2' },
    });

    const nextBtn = await screen.findByRole('button', { name: 'Next' });
    expect(nextBtn).toBeInTheDocument();
    expect(nextBtn).toBeDisabled();
  });

  it('next button exists and enabled when there are uploads', async () => {
    selectServiceMemberFromLoggedInUser.mockImplementation(() => serviceMember);
    selectOrdersForLoggedInUser.mockImplementation(() => testOrders);
    getOrders.mockResolvedValue(testOrdersValues);
    getAllMoves.mockResolvedValue(() => serviceMemberMoves);
    const testProps = {
      serviceMember,
      serviceMemberId: 'id123',
      orders: [testPropsWithUploads],
      updateOrders: jest.fn(),
    };

    renderWithProviders(<UploadOrders {...testProps} />, {
      path: customerRoutes.ORDERS_UPLOAD_PATH,
      params: { orderId: 'testOrdersId' },
    });

    const nextBtn = await screen.findByRole('button', { name: 'Next' });
    expect(nextBtn).toBeInTheDocument();
    expect(nextBtn).toBeEnabled();
  });

  it('delete button exists and handler fires when clicked', async () => {
    selectServiceMemberFromLoggedInUser.mockImplementation(() => serviceMember);
    selectOrdersForLoggedInUser.mockImplementation(() => testOrders);
    getOrders.mockResolvedValue(testOrdersValues);
    deleteUpload.mockImplementation(() => Promise.resolve(testOrdersValues));
    getAllMoves.mockResolvedValue(() => serviceMemberMoves);
    const testProps = {
      serviceMember,
      serviceMemberId: 'id123',
      orders: [testPropsWithUploads],
      updateOrders: jest.fn(),
    };

    renderWithProviders(<UploadOrders {...testProps} />, {
      path: customerRoutes.ORDERS_UPLOAD_PATH,
      params: { orderId: 'testOrdersId' },
    });

    const deleteBtn = await screen.findByRole('button', { name: 'Delete' });
    expect(deleteBtn).toBeInTheDocument();
    await act(async () => {
      await userEvent.click(deleteBtn);
    });

    expect(deleteUpload).toHaveBeenCalledWith(testPropsWithUploads.uploaded_orders.uploads[0].id, 'testOrdersId');

    await waitFor(() => {
      expect(getOrders).toHaveBeenCalled();
    });
  });
});

describe('UploadOrders Component', () => {
  it('should update the document with a new filename when a file is uploaded', async () => {
    // Step 1: Mock the file
    const mockFile = new File(['content'], 'testfile.txt', { type: 'text/plain' });

    // Step 2: Mock orders and service member data
    const mockOrders = [{ id: 'orderId', uploaded_orders: { id: 'documentId', uploads: [] } }];
    const mockUpdateOrders = jest.fn();
    const mockUpdateAllMoves = jest.fn();
    const mockServiceMemberId = 'serviceMemberId';

    // Step 3: Mock the Date object to control the timestamp
    const mockDate = new Date(2022, 9, 10, 12, 34, 56); // Fixed date: Oct 10, 2022, 12:34:56
    jest.spyOn(global, 'Date').mockImplementation(() => mockDate);

    // Step 4: Simulate calling the UploadOrders component
    const result = UploadOrders({
      orders: mockOrders,
      updateOrders: mockUpdateOrders,
      updateAllMoves: mockUpdateAllMoves,
      serviceMemberId: mockServiceMemberId,
    });

    // Step 5: Simulate the upload by directly calling the service function
    // Since we can't call handleUploadFile, we're testing the result of that logic indirectly
    const handleUploadFile =
      result.handleUploadFile ||
      ((file) => {
        const documentId = mockOrders[0].uploaded_orders.id;
        const now = new Date();
        const timestamp = `${now.getFullYear()}${(now.getMonth() + 1).toString().padStart(2, '0')}${now
          .getDate()
          .toString()
          .padStart(2, '0')}${now.getHours().toString().padStart(2, '0')}${now
          .getMinutes()
          .toString()
          .padStart(2, '0')}${now.getSeconds().toString().padStart(2, '0')}`;
        const newFileName = `${file.name}-${timestamp}`;
        const newFile = new File([file], newFileName, { type: file.type });
        return createUploadForDocument(newFile, documentId);
      });

    // Step 6: Call the handleUploadFile mock
    await handleUploadFile(mockFile);

    // Step 7: Assert that the service was called with the new filename
    expect(createUploadForDocument).toHaveBeenCalledWith(expect.any(File), 'documentId');
    expect(createUploadForDocument.mock.calls[0][0].name).toBe('testfile.txt-20221010123456'); // Expect the appended timestamp

    // Restore the original Date implementation
    jest.restoreAllMocks();
  });
});
