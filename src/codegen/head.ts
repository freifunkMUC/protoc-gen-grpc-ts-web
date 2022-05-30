import ts from 'typescript';

export function codegenHead(): ts.Statement[] {
  return [
    comment('Generated by protoc-gen-grpc-ts-web. DO NOT EDIT!'),
    comment('prettier-ignore'),
    comment('eslint-disable'),
    comment('tslint:disable'),
  ];
}

function comment(comment: string) {
  return ts.addSyntheticLeadingComment(
    ts.factory.createEmptyStatement(),
    ts.SyntaxKind.MultiLineCommentTrivia,
    ` ${comment} `,
    false
  );
}
