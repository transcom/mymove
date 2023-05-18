import React from 'react';
import { ShowController, Show, SimpleShowLayout, TextField, DateField } from 'react-admin';

const UploadShow = (props) => {
  return (
    <ShowController {...props}>
      {(controllerProps) => {
        const { record } = controllerProps;
        return (
          <Show {...controllerProps} {...record}>
            <SimpleShowLayout>
              {record && record.serviceMemberId
                ? [
                    <TextField key="serviceMemberId" source="serviceMemberId" label="Service Member ID" />,
                    <TextField
                      key="serviceMemberFirstName"
                      source="serviceMemberFirstName"
                      label="Service Member First Name"
                    />,
                    <TextField
                      key="serviceMemberLastName"
                      source="serviceMemberLastName"
                      label="Service Member Last Name"
                    />,
                    <TextField key="serviceMemberPhone" source="serviceMemberPhone" label="Service Member Phone" />,
                    <TextField key="serviceMemberEmail" source="serviceMemberEmail" label="Service Member Email" />,
                  ]
                : [
                    <TextField key="officeUserId" source="officeUserId" label="Office User ID" />,
                    <TextField key="officeUserFirstName" source="officeUserFirstName" label="Office User First Name" />,
                    <TextField key="officeUserLastName" source="officeUserLastName" label="Office User Last Name" />,
                    <TextField key="officeUserPhone" source="officeUserPhone" label="Office User Phone" />,
                    <TextField key="officeUserEmail" source="officeUserEmail" label="Office User Email" />,
                  ]}
              <TextField key="moveLocator" source="moveLocator" label="Move Locator" />
              <TextField key="upload.filename" source="upload.filename" label="Upload Filename" />
              <TextField key="upload.size" source="upload.size" label="Upload Size" />
              <TextField key="upload.contentType" source="upload.contentType" label="Upload Content Type" />
              <DateField key="upload.createdAt" source="upload.createdAt" showTime label="Created At" />
            </SimpleShowLayout>
          </Show>
        );
      }}
    </ShowController>
  );
};

export default UploadShow;
