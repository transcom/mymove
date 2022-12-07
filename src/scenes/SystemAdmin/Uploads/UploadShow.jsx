import React from 'react';
import { ShowController, ShowView, SimpleShowLayout, TextField, DateField, useRecordContext } from 'react-admin';

const record = useRecordContext();
const UploadShow = (props) => (
  <ShowController {...props}>
    {(record) => (
      <ShowView {...props} {...record}>
        <SimpleShowLayout>
          {record && record.serviceMemberId
            ? [
                <TextField source="serviceMemberId" label="Service Member ID" />,
                <TextField source="serviceMemberFirstName" label="Service Member First Name" />,
                <TextField source="serviceMemberLastName" label="Service Member Last Name" />,
                <TextField source="serviceMemberPhone" label="Service Member Phone" />,
                <TextField source="serviceMemberEmail" label="Service Member Email" />,
              ]
            : [
                <TextField source="officeUserId" label="Office User ID" />,
                <TextField source="officeUserFirstName" label="Office User First Name" />,
                <TextField source="officeUserLastName" label="Office User Last Name" />,
                <TextField source="officeUserPhone" label="Office User Phone" />,
                <TextField source="officeUserEmail" label="Office User Email" />,
              ]}
          <TextField source="moveLocator" label="Move Locator" />
          <TextField source="upload.filename" label="Upload Filename" />
          <TextField source="upload.size" label="Upload Size" />
          <TextField source="upload.contentType" label="Upload Content Type" />
          <DateField source="upload.createdAt" showTime label="Created At" />
        </SimpleShowLayout>
      </ShowView>
    )}
  </ShowController>
);

export default UploadShow;
