# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: shift.proto

from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='shift.proto',
  package='',
  syntax='proto3',
  serialized_options=None,
  serialized_pb=b'\n\x0bshift.proto\"\x80\x01\n\x06\x43mdMsg\x12\x1c\n\x03\x63md\x18\x01 \x01(\x0e\x32\x0f.CmdMsg.Command\x12\x0c\n\x04\x64\x61ta\x18\x02 \x01(\x0c\"J\n\x07\x43ommand\x12\x0f\n\x0bHealthCheck\x10\x00\x12\n\n\x06SetPin\x10\x01\x12\x0b\n\x07SetByte\x10\x02\x12\t\n\x05Reset\x10\x03\x12\n\n\x06Status\x10\x04\x62\x06proto3'
)



_CMDMSG_COMMAND = _descriptor.EnumDescriptor(
  name='Command',
  full_name='CmdMsg.Command',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='HealthCheck', index=0, number=0,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='SetPin', index=1, number=1,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='SetByte', index=2, number=2,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='Reset', index=3, number=3,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='Status', index=4, number=4,
      serialized_options=None,
      type=None),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=70,
  serialized_end=144,
)
_sym_db.RegisterEnumDescriptor(_CMDMSG_COMMAND)


_CMDMSG = _descriptor.Descriptor(
  name='CmdMsg',
  full_name='CmdMsg',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='cmd', full_name='CmdMsg.cmd', index=0,
      number=1, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='data', full_name='CmdMsg.data', index=1,
      number=2, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=b"",
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
    _CMDMSG_COMMAND,
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=16,
  serialized_end=144,
)

_CMDMSG.fields_by_name['cmd'].enum_type = _CMDMSG_COMMAND
_CMDMSG_COMMAND.containing_type = _CMDMSG
DESCRIPTOR.message_types_by_name['CmdMsg'] = _CMDMSG
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

CmdMsg = _reflection.GeneratedProtocolMessageType('CmdMsg', (_message.Message,), {
  'DESCRIPTOR' : _CMDMSG,
  '__module__' : 'shift_pb2'
  # @@protoc_insertion_point(class_scope:CmdMsg)
  })
_sym_db.RegisterMessage(CmdMsg)


# @@protoc_insertion_point(module_scope)
